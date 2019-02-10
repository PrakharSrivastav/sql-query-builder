/*
Package querybuilder helps to generate sql queries in different dialects.
This package can be best used with the scenarios where the structure of the domains models is unknown beforehand.

*/
package querybuilder

// QueryBuilder provides a contract to be implemented by all sql generators
type QueryBuilder interface {
	Create(columns map[string]interface{}) (string, error)
	Get(columns []string, where map[string]interface{}, limit, offset int) (string, error)
	Update(columns, where map[string]interface{}) (string, error)
	Insert(columns []string, data []map[string]interface{}) (string, error)
}
