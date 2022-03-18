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
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIEBDCCAuygAwIBAgIDAjppMA0GCSqGSIb3DQEBBQUAMEIxCzAJBgNVBAYTAlVT
MRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQDExJHZW9UcnVzdCBHbG9i
YWwgQ0EwHhcNMTMwNDA1MTUxNTU1WhcNMTUwNDA0MTUxNTU1WjBJMQswCQYDVQQG
EwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xlIEludGVy
bmV0IEF1dGhvcml0eSBHMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJwqBHdc2FCROgajguDYUEi8iT/xGXAaiEZ+4I/F8YnOIe5a/mENtzJEiaB0C1NP
VaTOgmKV7utZX8bhBYASxF6UP7xbSDj0U/ck5vuR6RXEz/RTDfRK/J9U3n2+oGtv
h8DQUB8oMANA2ghzUWx//zo8pzcGjr1LEQTrfSTe5vn8MXH7lNVg8y5Kr0LSy+rE
ahqyzFPdFUuLH8gZYR/Nnag+YyuENWllhMgZxUYi+FOVvuOAShDGKuy6lyARxzmZ
EASg8GF6lSWMTlJ14rbtCMoU/M4iarNOz0YDl5cDfsCx3nuvRTPPuj5xt970JSXC
DTWJnZ37DhF5iR43xa+OcmkCAwEAAaOB+zCB+DAfBgNVHSMEGDAWgBTAephojYn7
qwVkDBF9qn1luMrMTjAdBgNVHQ4EFgQUSt0GFhu89mi1dvWBtrtiGrpagS8wEgYD
VR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAQYwOgYDVR0fBDMwMTAvoC2g
K4YpaHR0cDovL2NybC5nZW90cnVzdC5jb20vY3Jscy9ndGdsb2JhbC5jcmwwPQYI
KwYBBQUHAQEEMTAvMC0GCCsGAQUFBzABhiFodHRwOi8vZ3RnbG9iYWwtb2NzcC5n
ZW90cnVzdC5jb20wFwYDVR0gBBAwDjAMBgorBgEEAdZ5AgUBMA0GCSqGSIb3DQEB
BQUAA4IBAQA21waAESetKhSbOHezI6B1WLuxfoNCunLaHtiONgaX4PCVOzf9G0JY
/iLIa704XtE7JW4S615ndkZAkNoUyHgN7ZVm2o6Gb4ChulYylYbc3GrKBIxbf/a/
zG+FA1jDaFETzf3I93k9mTXwVqO94FntT0QJo544evZG0R0SnU++0ED8Vf4GXjza
HFa9llF7b1cq26KqltyMdMKVvvBulRP/F/A8rLIQjcxz++iPAsbw+zOzlTvjwsto
WHPbqCRiOwY1nQ2pM714A5AuTHhdUDqB1O6gyHA43LL5Z/qHQF1hwFGPa4NrzQU6
yuGnBXj8ytqU0CwIPX4WecigUCAkVDNx
-----END CERTIFICATE-----`

var rsaCertPEM = `-----BEGIN CERTIFICATE-----
MIIB0zCCAX2gAwIBAgIJAI/M7BYjwB+uMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTIwOTEyMjE1MjAyWhcNMTUwOTEyMjE1MjAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJ
hPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wok/4xIA+ui35/MmNa
rtNuC+BdZ1tMuVCPFZcCAwEAAaNQME4wHQYDVR0OBBYEFJvKs8RfJaXTH08W+SGv
zQyKn0H8MB8GA1UdIwQYMBaAFJvKs8RfJaXTH08W+SGvzQyKn0H8MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADQQBJlffJHybjDGxRMqaRmDhX0+6v02TUKZsW
r5QuVbpQhH6u+0UgcW0jp9QwpxoPTLTWGXEWBBBurxFwiCBhkQ+V
-----END CERTIFICATE-----
`

var rsaKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBANLJhPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wo
k/4xIA+ui35/MmNartNuC+BdZ1tMuVCPFZcCAwEAAQJAEJ2N+zsR0Xn8/Q6twa4G
6OB1M1WO+k+ztnX/1SvNeWu8D6GImtupLTYgjZcHufykj09jiHmjHx8u8ZZB/o1N
MQIhAPW+eyZo7ay3lMz1V01WVjNKK9QSn1MJlb06h/LuYv9FAiEA25WPedKgVyCW
SmUwbPw8fnTcpqDWE3yTO3vKcebqMSsCIBF3UmVue8YU3jybC3NxuXq3wNm34R8T
xVLHwDXh/6NJAiEAl2oHGGLz64BuAfjKrqwz7qMYr9HCLIe/YsoWq/olzScCIQDi
D2lWusoe2/nEqfDVVWGWlyJ7yOmqaVm/iNUN9B2N2g==
-----END RSA PRIVATE KEY-----
`

func TestOpenTLS(t *testing.T) {
	type args struct {
		tlsName             string
		mysqlCaCertPath     string
		mysqlClientCertPath string
		mysqlClientKeyPath  string
	}
	type file struct {
		mysqlCaCertValue     string
		mysqlClientCertValue string
		mysqlClientKeyValue  string
	}
	tests := []struct {
		name    string
		args    args
		file    file
		wantErr bool
	}{
		{
			name:    "tlsName was empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "mysqlCaCertPath was empty",
			args: args{
				tlsName: "tlsName",
			},
			wantErr: false,
		},
		{
			name: "test ca cert",
			args: args{
				tlsName:         "tlsName",
				mysqlCaCertPath: "true",
			},
			file: file{
				mysqlCaCertValue: rootPEM,
			},
			wantErr: false,
		},
		{
			name: "test ca cert and client cert",
			args: args{
				tlsName:             "tlsName",
				mysqlCaCertPath:     "true",
				mysqlClientCertPath: "true",
				mysqlClientKeyPath:  "true",
			},
			file: file{
				mysqlCaCertValue:     rootPEM,
				mysqlClientKeyValue:  rsaKeyPEM,
				mysqlClientCertValue: rsaCertPEM,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errorInfo error
			if tt.args.mysqlCaCertPath != "" {
				f, err := os.CreateTemp("", "ca-cert")
				if err != nil {
					errorInfo = err
				} else {
					_, err = f.Write([]byte(tt.file.mysqlCaCertValue))
					if err != nil {
						errorInfo = err
					}
					tt.args.mysqlCaCertPath = f.Name()
				}
			}

			if tt.args.mysqlClientCertPath != "" && tt.args.mysqlClientKeyPath != "" {
				certFile, err := os.CreateTemp("", "client-cert")
				if err != nil {
					errorInfo = err
				} else {
					_, err = certFile.Write([]byte(tt.file.mysqlClientCertValue))
					if err != nil {
						errorInfo = err
					}
					tt.args.mysqlClientCertPath = certFile.Name()
				}
				keyFile, err := os.CreateTemp("", "client-key")
				if err != nil {
					errorInfo = err
				} else {
					_, err = keyFile.Write([]byte(tt.file.mysqlClientKeyValue))
					if err != nil {
						errorInfo = err
					}
					tt.args.mysqlClientKeyPath = keyFile.Name()
				}
			}

			if errorInfo != nil {
				tt.wantErr = true
			}

			patch := monkey.Patch(mysql.RegisterTLSConfig, func(key string, config *tls.Config) error {
				if tt.args.tlsName != "" && tt.args.mysqlCaCertPath != "" {
					assert.NotNil(t, key)
					assert.Equal(t, tt.args.tlsName, key)
					assert.NotNil(t, config.RootCAs)
				}
				if tt.args.mysqlClientKeyPath != "" && tt.args.mysqlClientCertPath != "" {
					assert.NotNil(t, config.Certificates)
					assert.Equal(t, 1, len(config.Certificates))
				}
				return nil
			})
			defer patch.Unpatch()

			if err := OpenTLS(tt.args.tlsName, tt.args.mysqlCaCertPath, tt.args.mysqlClientCertPath, tt.args.mysqlClientKeyPath); (err != nil) != tt.wantErr {
				t.Errorf("OpenTLS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
