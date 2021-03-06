package shared

import (
	"bytes"
	"crypto/x509"
	"testing"

	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp/constants"
	keystore "github.com/pavel-v-chernykh/keystore-go"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestEnvOverride(t *testing.T) {
	src := []corev1.EnvVar{
		{
			Name:  "test1",
			Value: "value1",
		},
		{
			Name:  "test2",
			Value: "value2",
		},
	}
	dst := []corev1.EnvVar{
		{
			Name:  "test1",
			Value: "valueX",
		},
		{
			Name:  "test3",
			Value: "value3",
		},
	}
	result := EnvOverride(src, dst)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, result[0], dst[0])
	assert.Equal(t, result[1], src[1])
	assert.Equal(t, result[2], dst[1])
}

func TestGetEnvVar(t *testing.T) {
	vars := []corev1.EnvVar{
		{
			Name:  "test1",
			Value: "value1",
		},
		{
			Name:  "test2",
			Value: "value2",
		},
	}
	pos := GetEnvVar("test1", vars)
	assert.Equal(t, 0, pos)

	pos = GetEnvVar("other", vars)
	assert.Equal(t, -1, pos)
}

func TestGenerateKeystore(t *testing.T) {
	password := GeneratePassword(8)
	assert.EqualValues(t, 8, len(password))

	commonName := "test-https"
	keyBytes := GenerateKeystore(commonName, password)
	keyStore, err := keystore.Decode(bytes.NewReader(keyBytes), password)
	assert.Nil(t, err)

	derKey := keyStore[constants.KeystoreAlias].(*keystore.PrivateKeyEntry).PrivKey
	_, err = x509.ParsePKCS8PrivateKey(derKey)
	assert.Nil(t, err)

	cert := keyStore[constants.KeystoreAlias].(*keystore.PrivateKeyEntry).CertChain[0].Content
	certificate, err := x509.ParseCertificate(cert)
	assert.Nil(t, err)
	assert.Equal(t, commonName, certificate.Subject.CommonName)
}

func TestEnvVarCheck(t *testing.T) {
	empty := []corev1.EnvVar{}
	a := []corev1.EnvVar{
		{Name: "A", Value: "1"},
	}
	b := []corev1.EnvVar{
		{Name: "A", Value: "2"},
	}
	c := []corev1.EnvVar{
		{Name: "A", Value: "1"},
		{Name: "B", Value: "1"},
	}

	assert.True(t, EnvVarCheck(empty, empty))
	assert.True(t, EnvVarCheck(a, a))

	assert.False(t, EnvVarCheck(empty, a))
	assert.False(t, EnvVarCheck(a, empty))

	assert.False(t, EnvVarCheck(a, b))
	assert.False(t, EnvVarCheck(b, a))

	assert.False(t, EnvVarCheck(a, c))
	assert.False(t, EnvVarCheck(c, b))
}
