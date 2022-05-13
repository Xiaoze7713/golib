/**
 * @Author: wenliangzhang
 * @Description:
 * @File: mysql
 * @Version: 1.0.0
 * @Date: 2022/5/9 8:26 PM
 */
package mysql

func GetMysqlInstence(dbName string) GormConnector {
	return dbConnectors[dbName]
}
