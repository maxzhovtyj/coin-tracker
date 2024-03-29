package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/service"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	"math"
	"strconv"
	"time"
)

const (
	subscribeCoinNameInputStep     = "subscription coin name"
	subscribeCoinIntervalInputStep = "subscription notify interval"
)

func (h *Handler) Subscriptions(ctx *Context) {
	subscriptions, err := h.service.Subscription.UserSubscriptions(ctx.UID)
	if err != nil {
		ctx.ResponseString("Can't get your subscriptions, " + err.Error())
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, "Here is the list of your subscriptions")
	btns := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("New Coin Subscription", subscribeCoinCallback+"="+"")},
	}
	btns = append(btns, h.subscriptionsInlineKeyboard(subscriptions)...)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(btns...)
	ctx.Response(msg)
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

func (h *Handler) SubscribeCoin(ctx *Context) {
	ctx.ResponseString("Please enter coin name to subscribe")
	ctx.FSM.Update(ctx.UID, State{
		Caller: ctx.CallbackName,
		Step:   subscribeCoinNameInputStep,
	})
}

func (h *Handler) subscribeCoinNameInputStep(ctx *Context) {
	coinName := ctx.Update.Message.Text

	ctx.ResponseString("Please enter notifying interval (Ex.: 30m, 3h, 1d)")

	ctx.FSM.Update(ctx.UID, State{
		Caller: ctx.FSM.Get(ctx.UID).Caller,
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

func (h *Handler) SubscriptionInfo(ctx *Context) {
	id, err := strconv.ParseInt(ctx.CallbackDataValue, 10, 64)
	if err != nil {
		ctx.ResponseString("Invalid subscription id")
		return
	}

	s, err := h.service.Subscription.Get(id)
	if err != nil {
		ctx.ResponseString("Can't get subscription, " + err.Error())
		return
	}

	msg := tgbotapi.NewMessage(ctx.UID, formatSubscription(s))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Delete", fmt.Sprintf("%s=%d", subscriptionDeleteCallback, s.ID)),
	})

	ctx.Response(msg)
}

func formatSubscription(s models.Subscription) string {
	var d string

	switch s.Type {
	case service.CoinSubscriptionType:
		var cs models.CoinSubscriptionData

		if err := json.Unmarshal([]byte(s.Data), &cs); err != nil {
			return ""
		}

		d = fmt.Sprintf(`Coin Name: %s`, cs.CoinName)
	}

	return fmt.Sprintf(`
Type: %s
Interval: %s
%s
Last notified at: %s
`, s.Type, s.NotifyInterval, d, s.LastNotifiedAt.Format(time.DateTime))
}

func (h *Handler) DeleteSubscription(ctx *Context) {
	sid, err := strconv.ParseInt(ctx.CallbackDataValue, 10, 64)
	if err != nil {
		ctx.ResponseString("Can't find subscription")
		return
	}

	err = h.service.Subscription.Delete(sid)
	if err != nil {
		ctx.ResponseString("Can't delete subscription")
		return
	}

	ctx.ResponseString("Subscription successfully deleted")
}

func (h *Handler) SubscriptionsWorker() {
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

func (h *Handler) subscriptionsInlineKeyboard(subs []models.Subscription) [][]tgbotapi.InlineKeyboardButton {
	rows := math.Ceil(float64(len(subs)) / float64(2))
	keyboard := make([][]tgbotapi.InlineKeyboardButton, int(rows))

	var row int
	var col int

	for _, t := range subs {
		var d models.CoinSubscriptionData

		err := json.Unmarshal([]byte(t.Data), &d)
		if err != nil {
			h.logger.Error(err)
			continue
		}

		cbData := fmt.Sprintf("%s=%d", subscriptionCoinInfoCallback, t.ID)
		keyboard[row] = append(keyboard[row], tgbotapi.NewInlineKeyboardButtonData(d.CoinName, cbData))

		if (col+1)%2 == 0 {
			col = 0
			row++
		} else {
			col++
		}
	}

	return keyboard
}
