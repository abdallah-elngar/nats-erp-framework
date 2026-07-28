package signals

import (
    "log"

    "github.com/nats-framework/nats/apps/sale/models"
    "github.com/nats-framework/nats/pkg/signals"
)

// RegisterSignals يسجل الإشارات
func RegisterSignals(signalManager *signals.SignalManager) {
    // إشارة عند تغيير Sale
    signalManager.Connect("sale.sale.changed", func(signal *signals.Signal) error {
        item, ok := signal.Data.(*models.Sale)
        if !ok {
            return nil
        }
        log.Printf("📊 Sale changed: sa_name: %+v, price: %+v, quantity: %+v", item.SaName, item.Price, item.Quantity)
        return nil
    })

    // ✅ إشارة عند تغيير حالة معينة (مثال: انخفاض المخزون)
    signalManager.Connect("sale.sale.low_stock", func(signal *signals.Signal) error {
        item, ok := signal.Data.(*models.Sale)
        if !ok {
            return nil
        }
        log.Printf("⚠️ Low stock alert for Sale: sa_name: %+v, price: %+v, quantity: %+v", item.SaName, item.Price, item.Quantity)
        return nil
    })
}
