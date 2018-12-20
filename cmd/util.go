package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/pgpkeys"
)

func formatFingerprint(fingerprint string) string {
	return fmt.Sprintf("%s %s %s %s", strings.ToUpper(fingerprint[24:])[0:4],
		strings.ToUpper(fingerprint[24:])[4:8],
		strings.ToUpper(fingerprint[24:])[8:12],
		strings.ToUpper(fingerprint[24:])[12:16])
}

func keybaseEncrypt(keybaseIdentity string, payload []byte) (string, []byte, error) {
	plaintextes := make([][]byte, 0)
	plaintextes = append(plaintextes, payload)
	pgpKeys := make([]string, 0)
	pgpKeys = append(pgpKeys, keybaseIdentity)
	pgpKeysFetched, err := pgpkeys.FetchKeybasePubkeys(pgpKeys)
	if err != nil {
		return "", nil, err
	}
	keys := make([]string, 0)
	for _, fetched := range pgpKeysFetched {
		keys = append(keys, fetched)
	}
	fingerprints, ciphertextes, err := pgpkeys.EncryptShares(plaintextes, keys)
	if err != nil {
		return "", nil, err
	}
	return formatFingerprint(fingerprints[0]), ciphertextes[0], nil
}
