/**
 * 数据库操作
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/1 11:08
 */

package duang

import (
	"boom/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"reflect"
	"time"
)

var (
	tablePrefix = config.RetDatabase().Prefix
	db          *gorm.DB
)

func init() {
	db = RetDb()
}

// RetDb 初始化数据库连接
func RetDb() *gorm.DB {
	db, _ = gorm.Open(mysql.Open(config.RetDatabase().Dns), &gorm.Config{})
	sqlDB, er := db.DB()
	if er != nil {
		log.Println(er)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute)
	return db
}

func Pdo_get(tableName string, columns []string, where map[string]interface{}) map[string]interface{} {
	var results map[string]interface{}
	var order string
	for k, v := range where {
		if k == "Order" {
			order = v.(string)
			delete(where, k)
		}
	}
	if order != "" {
		db.Table(tablePrefix+tableName).Select(columns).Order(order).Find(&results, where)
	} else {
		db.Table(tablePrefix+tableName).Select(columns).Find(&results, where)
	}
	return results
}

func Pdo_getall(tableName string, columns []string, where map[string]interface{}) (result []map[string]interface{}, pageNum int64, total int64) {
	var results []map[string]interface{}
	var num int      // 显示条数
	var order string // 排序
	var page int     // 页码
	var pageSize int // 分页大小
	for k, v := range where {
		if k == "Limit" {
			if reflect.TypeOf(v).Kind() == reflect.Int { // 判断是否要分页
				num = v.(int)
			} else {
				page = v.([2]int)[0]
				pageSize = v.([2]int)[1]
				if page <= 0 {
					page = 1
				}
			}
			delete(where, k)
		} else if k == "Order" {
			order = v.(string)
			delete(where, k)
		}
	}
	if page != 0 && pageSize != 0 {
		db.Table(tablePrefix + tableName).Where(where).Count(&total)
		pageNum = total / int64(pageSize)
		if total%int64(pageSize) != 0 {
			pageNum++
		}
	}
	if num != 0 || page != 0 && order != "" {
		if page != 0 && pageSize != 0 {
			db.Table(tablePrefix+tableName).Select(columns).Order(order).Offset((page-1)*pageSize).Limit(pageSize).Find(&results, where)
			//[]map[string]interface{}{{"pageNum": pageNum, "total": total,"results":results}}
			return results, pageNum, total
		} else {
			db.Table(tablePrefix+tableName).Select(columns).Order(order).Limit(num).Find(&results, where)
		}
	} else if num != 0 || page != 0 {
		if page != 0 && pageSize != 0 {
			db.Table(tablePrefix+tableName).Select(columns).Limit(pageSize).Offset((page-1)*pageSize).Find(&results, where)
			//[]map[string]interface{}{{"pageNum": pageNum, "total": total,"results":results}}
			return results, pageNum, total
		} else {
			db.Table(tablePrefix+tableName).Select(columns).Limit(num).Find(&results, where)
		}
	} else if order != "" {
		db.Table(tablePrefix+tableName).Select(columns).Order(order).Find(&results, where)
	} else {
		db.Table(tablePrefix+tableName).Select(columns).Find(&results, where)
	}
	return results, pageNum, total
}

func Pdo_insert[T map[string]interface{} | []map[string]interface{}](tableName string, data T) int64 {
	result := db.Table(tablePrefix + tableName).Create(&data)
	if result.Error != nil {
		log.Panicln(result.Error)
	}
	return result.RowsAffected
}

func Pdo_delete(tableName string, where map[string]interface{}) int64 {
	var results map[string]interface{}
	res := db.Table(tablePrefix+tableName).Delete(&results, where)
	return res.RowsAffected
}

func Pdo_update(tableName string, data map[string]interface{}, where map[string]interface{}) int64 {
	res := db.Table(tablePrefix + tableName).Where(where).Updates(data)
	if res.Error != nil {
		log.Panicln(res.Error)
	}
	return res.RowsAffected
}

func Pdo_count(tableName string, where map[string]interface{}) int64 {
	var total int64
	res := db.Table(tablePrefix + tableName).Where(where).Count(&total)
	if res.Error != nil {
		log.Panicln(res.Error)
	}
	return total
}
