# 🚀 NATS Framework

**NATS** هو إطار عمل متكامل (Full-Stack Framework) مبني بلغة Go، يهدف إلى تبسيط وتسريع تطوير أنظمة المؤسسات (ERP) وأنظمة إدارة الأعمال.

## ✨ المميزات

- **Zero Code**: المطور لا يكتب أي كود يدوياً
- **CLI Interactive**: أوامر تفاعلية مع أسئلة وأمثلة
- **Auto Generation**: كل الملفات تنشأ تلقائياً
- **Dual Migrations**: GORM + SQL معاً
- **HTMX Frontend**: واجهات تفاعلية بدون JavaScript ثقيل
- **Local Static**: جميع الملفات الثابتة محلية
- **RBAC**: نظام صلاحيات متكامل
- **Multi-Database**: دعم PostgreSQL, MySQL, SQLite

## 🚀 البدء السريع

```bash
# تنزيل النظام
git clone https://github.com/nats-framework/nats.git
cd nats

# تثبيت الاعتماديات
go mod download

# بناء النظام
make build

# إنشاء مشروع جديد
nats init my-erp-system

# تشغيل الخادم
cd my-erp-system
nats serve