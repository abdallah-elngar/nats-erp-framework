package hooks

import (
    "log"

    "github.com/nats-framework/nats/apps/sale/models"
    "github.com/nats-framework/nats/pkg/hooks"
)

// RegisterHooks يسجل الخطافات
func RegisterHooks(hookManager *hooks.HookManager) {
    // خطاف قبل إنشاء Sale
    hookManager.RegisterBefore("sale.sale.create", func(data interface{}) (interface{}, error) {
        item, ok := data.(*models.Sale)
        if !ok {
            return data, nil
        }
        log.Printf("🔧 Before creating Sale: sa_name: %+v, price: %+v, quantity: %+v", item.SaName, item.Price, item.Quantity)
        // ✅ يمكن إضافة تحقق أو تعديل قبل الإنشاء
        // if item.Price < 0 {
        //     return nil, fmt.Errorf("price cannot be negative")
        // }
        return data, nil
    }, 10)

    // خطاف بعد إنشاء Sale
    hookManager.RegisterAfter("sale.sale.create", func(data interface{}) (interface{}, error) {
        item, ok := data.(*models.Sale)
        if !ok {
            return data, nil
        }
        log.Printf("✅ Sale created: sa_name: %+v, price: %+v, quantity: %+v", item.SaName, item.Price, item.Quantity)
        return data, nil
    }, 10)
}
