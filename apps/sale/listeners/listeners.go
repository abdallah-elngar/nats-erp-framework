package listeners

import (
    "fmt"
    "log"

    "github.com/nats-framework/nats/apps/sale/models"
    "github.com/nats-framework/nats/pkg/events"
)

// RegisterListeners يسجل مستمعي الأحداث
func RegisterListeners(em *events.EventManager) {
    // مستمع عند إنشاء Sale
    em.Listen("sale.sale.created", func(event events.Event) error {
        item, ok := event.GetData().(*models.Sale)
        if !ok {
            return fmt.Errorf("invalid event data")
        }
        log.Printf("📦 Sale created: sa_name: %+v, price: %+v, quantity: %+v", item.SaName, item.Price, item.Quantity)
        return nil
    })

    // مستمع عند تحديث Sale
    em.Listen("sale.sale.updated", func(event events.Event) error {
        item, ok := event.GetData().(*models.Sale)
        if !ok {
            return fmt.Errorf("invalid event data")
        }
        log.Printf("📦 Sale updated: sa_name: %+v, price: %+v, quantity: %+v", item.SaName, item.Price, item.Quantity)
        return nil
    })

    // مستمع عند حذف Sale
    em.Listen("sale.sale.deleted", func(event events.Event) error {
        item, ok := event.GetData().(*models.Sale)
        if !ok {
            return fmt.Errorf("invalid event data")
        }
        log.Printf("📦 Sale deleted: ID=%d", item.ID)
        return nil
    })
}
