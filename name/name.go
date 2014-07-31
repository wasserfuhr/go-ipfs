package naming

import (
    "net"
    "strings"
    "regexp"
    "errors"
    proquint "github.com/Bren2010/proquint"
    mh "github.com/jbenet/go-multihash"
)

func Resolve(name string) (mh.Multihash, error) {
    b58Exp := "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$"
    pqExp := "^([abdfghijklmnoprstuvz]{5}-)*[abdfghijklmnoprstuvz]{5}$"

    isB58, err := regexp.MatchString(b58Exp, name)
    isPQ, err := regexp.MatchString(pqExp, name)

    if err != nil {
        return nil, err
    }

    if isB58 { // Is a base58 hash.
        h, err := mh.FromB58String(name)
        return h, err
    } else if isPQ { // Is a Proquint identifier.
        return mh.Multihash(proquint.Decode(name)), nil
    } else { // Is a domain name.  Hopefully.
        txts, err := net.LookupTXT(name)
        if err != nil {
            return nil, err
        }

        for i := 0; i < len(txts); i++ {
            var parts []string = strings.SplitN(txts[i], "=", 2)

            if len(parts) == 2 && parts[0] == "ipfs" {
                h, err := mh.FromB58String(parts[1])
                return h, err
            }
        }

        return nil, errors.New("Could not resolve IPNS.")
    }
}
