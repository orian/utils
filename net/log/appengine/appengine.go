package appengine

import (
	"github.com/orian/utils/net/log"
	ae "google.golang.org/appengine/log"
)

func init() {
	log.Debugf = ae.Debugf
	log.Infof = ae.Infof
	log.Warningf = ae.Warningf
	log.Errorf = ae.Errorf
	log.Criticalf = ae.Criticalf
}
