/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
)

//
// This just establishes the network connection; once established the type of connection should be irrelevant
//

// openConnection opens a network connection using the protocol indicated in the URL.
// Does not carry out any MQTT specific handshakes.
func openConnection(uri *url.URL, tlsc *tls.Config, timeout time.Duration, headers http.Header, websocketOptions *WebsocketOptions) (net.Conn, error) {
	switch uri.Scheme {
	case "ws":
		conn, err := NewWebsocket(uri.String(), nil, timeout, headers, websocketOptions)
		return conn, err
	case "wss":
		conn, err := NewWebsocket(uri.String(), tlsc, timeout, headers, websocketOptions)
		return conn, err
	case "mqtt", "tcp":
		allProxy := os.Getenv("all_proxy")
		if len(allProxy) == 0 {
			// FIXME - Add top-level APIs to optionally provide the
			// net.Dialer or dialer control function
			d := net.Dialer{
				Timeout: timeout,
				Control: func(network, address string, c syscall.RawConn) error {
					return c.Control(func(fd uintptr) {
						// set the socket options
						JUNOS_IP_RTBL_INDEX := 109
						JUNOS_IRT_VAL := 1
						err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, JUNOS_IP_RTBL_INDEX, JUNOS_IRT_VAL)
						//err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, 1)
						if err != nil {
							log.Println("setsocketopt: ", err)
						}
					})
				},
			}
			conn, err := d.Dial("tcp", uri.Host)
			//conn, err := net.DialTimeout("tcp", uri.Host, timeout)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		proxyDialer := proxy.FromEnvironment()

		conn, err := proxyDialer.Dial("tcp", uri.Host)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case "unix":
		conn, err := net.DialTimeout("unix", uri.Host, timeout)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case "ssl", "tls", "mqtts", "mqtt+ssl", "tcps":
		allProxy := os.Getenv("all_proxy")
		if len(allProxy) == 0 {
			conn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", uri.Host, tlsc)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		proxyDialer := proxy.FromEnvironment()

		conn, err := proxyDialer.Dial("tcp", uri.Host)
		if err != nil {
			return nil, err
		}

		tlsConn := tls.Client(conn, tlsc)

		err = tlsConn.Handshake()
		if err != nil {
			_ = conn.Close()
			return nil, err
		}

		return tlsConn, nil
	}
	return nil, errors.New("unknown protocol")
}
