package mysql

import (
	"fmt"
	"github.com/saperliu/common-tool/vo"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	mysqlClient := Database{}
	mysqlClient.User = "root"
	mysqlClient.Password = "enerbos,123"
	mysqlClient.Host = "127.0.0.1"
	mysqlClient.Port = 3306
	mysqlClient.TablePrefix = "point"
	mysqlClient.Database = "cron_db"
	mysqlClient.MaxOpenConns = 2
	mysqlClient.MaxIdleConns = 1
	_, err := mysqlClient.CreateDB()
	fmt.Printf("  %s ", err)
	store := vo.StorePoint{}
	store.Status = 1
	store.DataType = "normal"
	store.Tagid = "66852"
	store.Value = "223.24"
	store.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	mysqlClient.SaveStorePoint(store)
}
