// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysqldriver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// OpenTLS according to the passed parameters, decide whether to register tls config for mysql driver
func OpenTLS(tlsName, mysqlCaCertPath, mysqlClientCertPath, mysqlClientKeyPath string) error {
	if tlsName == "" || mysqlCaCertPath == "" {
		return nil
	}

	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(mysqlCaCertPath)
	if err != nil {
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append PEM")
	}

	if mysqlClientCertPath == "" || mysqlClientKeyPath == "" {
		// skip client cert
		err = mysql.RegisterTLSConfig(tlsName, &tls.Config{
			RootCAs: rootCertPool,
		})
		return err
	}

	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(mysqlClientCertPath, mysqlClientKeyPath)
	if err != nil {
		return fmt.Errorf("failed to append client PEM %v", err)
	}
	clientCert = append(clientCert, certs)
	// two-way encryption
	err = mysql.RegisterTLSConfig(tlsName, &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
	})
	return err
}
