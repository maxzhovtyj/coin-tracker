package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/service"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	"math"
	"time"
)

const (
	subscribeCoinNameInputStep     = "subscription coin name"
	subscribeCoinIntervalInputStep = "subscription notify interval"
)

func (h *Handler) SubscribeCoin(ctx *Context) {
	ctx.ResponseString("Please enter coin name to subscribe")
	ctx.FSM.Update(ctx.UID, State{
		Caller: subscribeCoinCommand.String(),
		Step:   subscribeCoinNameInputStep,
	})
}

type SubscribeCoinInput struct {
	CoinName string
	Interval time.Duration
}

func (h *Handler) ResolveSubscribeCoinSteps(ctx *Context) {
	state := ctx.FSM.Get(ctx.UID)

	switch state.Step {
	case subscribeCoinNameInputStep:
		h.subscribeCoinNameInputStep(ctx)
	case subscribeCoinIntervalInputStep:
		h.subscribeCoinIntervalInputStep(ctx)
	}
}

func (h *Handler) subscribeCoinNameInputStep(ctx *Context) {
	coinName := ctx.Update.Message.Text

	ctx.ResponseString("Please enter notifying interval (Ex.: 30m, 3h, 1d)")

	ctx.FSM.Update(ctx.UID, State{
		Caller: subscribeCoinCommand.String(),
		Step:   subscribeCoinIntervalInputStep,
		Data: SubscribeCoinInput{
			CoinName: coinName,
		},
	})
}

func (h *Handler) subscribeCoinIntervalInputStep(ctx *Context) {
	interval, err := time.ParseDuration(ctx.Update.Message.Text)
	if err != nil {
		ctx.ResponseString("Sorry, I can't parse this interval, try again (Ex.: 30m, 3h, 1d)")
		return
	}

	defer ctx.FSM.Remove(ctx.UID)

	data, ok := ctx.FSM.Get(ctx.UID).Data.(SubscribeCoinInput)
	if !ok {
		ctx.ResponseString("Sorry, can't get your input data")
		return
	}

	err = h.service.Subscription.NewCoinSubscription(ctx.UID, data.CoinName, interval)
	if err != nil {
		ctx.ResponseString("Sorry, I can't create new subscription, " + err.Error())
		return
	}

	ctx.ResponseString("New subscription created")
}

func (h *Handler) Subscriptions() {
	for range time.Tick(5 * time.Minute) {
		h.processAllSubscriptions()
	}
}

func (h *Handler) processAllSubscriptions() {
	start := time.Now()

	all, err := h.service.Subscription.All()
	if err != nil {
		h.logger.Errorf("can't retrieve subscriptions: %v", err)
		return
	}

	for _, s := range all {
		switch s.Type {
		case service.CoinSubscriptionType:
			h.processCoinSubscription(s)
		}
	}

	h.logger.Infof("finished processing subscriptions, took %s", time.Since(start))
}

func (h *Handler) processCoinSubscription(s models.Subscription) {
	var subData models.CoinSubscriptionData
	if err := json.Unmarshal([]byte(s.Data), &subData); err != nil {
		h.logger.Error(err)
		return
	}

	if time.Since(s.LastNotifiedAt) < s.NotifyInterval {
		return
	}

	ticker, err := h.service.Subscription.CoinTicker(subData.CoinName, subData.Interval)
	if err != nil {
		h.logger.Error(err)
		return
	}

	msg := tgbotapi.NewMessage(s.ChatID, h.coinSubscriptionText(ticker))

	_, err = h.bot.Send(msg)
	if err != nil {
		h.logger.Error(err)
		return
	}

	err = h.service.Subscription.Notified(s.ID)
	if err != nil {
		h.logger.Error(err)
		return
	}
}

func (h *Handler) coinSubscriptionText(ticker binance.SymbolTicker) string {
	return fmt.Sprintf(`
Here is your %s coin update 💲🤑💰:

  Price Change %s
  Price Change Percent %s
  Average Price %s
  Open Price %s
  High Price %s
  Low Price %s
  Last Price %s

-----------------------------------
`,
		ticker.Symbol,
		ticker.PriceChange,
		ticker.PriceChangePercent,
		ticker.WeightedAvgPrice,
		ticker.OpenPrice,
		ticker.HighPrice,
		ticker.LowPrice,
		ticker.LastPrice,
	)
}

func (h *Handler) CancelCoinSubscription(ctx *Context) {
	subscriptions, err := h.service.Subscription.UserSubscriptions(ctx.UID)
	if err != nil {
		ctx.ResponseString("Can't get your subscriptions, " + err.Error())
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, "Here is the list of your subscriptions")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(h.subscriptionsInlineKeyboard(subscriptions)...)
	ctx.Response(msg)
}

func (h *Handler) subscriptionsInlineKeyboard(subs []models.Subscription) [][]tgbotapi.InlineKeyboardButton {
	rows := math.Ceil(float64(len(subs)) / float64(2))
	keyboard := make([][]tgbotapi.InlineKeyboardButton, int(rows))

	var row int
	var col int

	for _, t := range subs {
		cbData := fmt.Sprintf("%s=%s", "cancelSubscription", t.Type)
		keyboard[row] = append(keyboard[row], tgbotapi.NewInlineKeyboardButtonData(t.Type, cbData))

		if (col+1)%2 == 0 {
			col = 0
			row++
		} else {
			col++
		}
	}

	return keyboard
}
