package inventory

import (
	"atlas-messages/character"
	"atlas-messages/command"
	"context"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func AwardItemCommandProducer(l logrus.FieldLogger) func(ctx context.Context) func(worldId byte, channelId byte, c character.Model, m string) (command.Executor, bool) {
	return func(ctx context.Context) func(worldId byte, channelId byte, c character.Model, m string) (command.Executor, bool) {
		return func(worldId byte, channelId byte, c character.Model, m string) (command.Executor, bool) {
			if !c.Gm() {
				l.Debugf("Ignoring character [%d] command [%s], because they are not a gm.", c.Id(), m)
				return nil, false
			}

			if !strings.HasPrefix(m, "@award item") {
				return nil, false
			}

			var itemId uint32
			var quantity uint16

			re := regexp.MustCompile("@award item (\\d*) (\\d*)")
			match := re.FindStringSubmatch(m)
			if len(match) == 3 {
				tItemId, err := strconv.ParseUint(match[1], 10, 32)
				if err != nil {
					return nil, false
				}
				itemId = uint32(tItemId)

				tQuantity, err := strconv.ParseInt(match[2], 10, 16)
				if err != nil {
					return nil, false
				}
				quantity = uint16(tQuantity)
			} else {
				re = regexp.MustCompile("@award item (\\d*)")
				match = re.FindStringSubmatch(m)
				if len(match) == 2 {
					tItemId, err := strconv.ParseUint(match[1], 10, 32)
					if err != nil {
						return nil, false
					}
					itemId = uint32(tItemId)
					quantity = 1
				} else {
					return nil, false
				}
			}

			exists := Exists(l)(ctx)(itemId)
			if !exists {
				l.Debugf("Ignoring character [%d] command [%s], because they did not input a valid item.", c.Id(), m)
				return nil, false
			}

			return func(l logrus.FieldLogger) func(ctx context.Context) error {
				return func(ctx context.Context) error {
					return GainItem(l, ctx)(c.Id(), itemId, quantity)
				}
			}, true
		}
	}
}

func GainItem(l logrus.FieldLogger, ctx context.Context) func(characterId uint32, itemId uint32, quantity uint16) error {
	return func(characterId uint32, itemId uint32, quantity uint16) error {
		_, err := CreateItem(l)(ctx)(characterId, itemId, quantity)
		return err
	}
}
