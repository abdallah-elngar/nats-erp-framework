#!/bin/bash

# ============================================
# نسخ احتياطي لقاعدة البيانات
# ============================================

set -e

echo "💾 Creating database backup..."

# المتغيرات
BACKUP_DIR="./storage/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.sql"

# إنشاء مجلد النسخ الاحتياطي
mkdir -p $BACKUP_DIR

# تحميل المتغيرات البيئية
source .env

# إنشاء النسخة الاحتياطية
if [ "$DB_DRIVER" = "postgres" ]; then
    PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME > $BACKUP_FILE
elif [ "$DB_DRIVER" = "mysql" ]; then
    mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME > $BACKUP_FILE
else
    echo "⚠️ Unsupported database driver: $DB_DRIVER"
    exit 1
fi

# ضغط الملف
gzip $BACKUP_FILE

echo "✅ Backup complete: $BACKUP_FILE.gz"
echo "📊 Size: $(du -h $BACKUP_FILE.gz | cut -f1)"

# حذف النسخ الاحتياطية القديمة (احتفظ بآخر 30)
find $BACKUP_DIR -name "backup_*.sql.gz" -type f -mtime +30 -delete
echo "🧹 Cleaned up old backups"