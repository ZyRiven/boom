/**
 * This file is part of Duang framework.
 *
 * Licensed under The MIT License
 * For full copyright and license information, please see the MIT-LICENSE.txt
 * Redistributions of files must retain the above copyright notice.
 *
 * @author    公子露<moralbodhi@gmail.com>
 * @copyright N+ LAB<TIMESCENTER>
 * @link      http://cloud.timescenter.ncn/
 * @license   http://www.opensource.org/licenses/mit-license.php MIT License
 */

package duang

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gohello/config"
	"net/http"
	"os"
	"strings"
)

// HandlerFunc defines the request handler used by duang
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // store all groups
	}
)

// New is the constructor of duang.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	if len(prefix) > 0 && prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix == "/" {
		prefix = ""
	}
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	//fmt.Println(method + pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	AscllTable()
	return server.ListenAndServe()
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func AscllTable() {
	fmt.Println("      .______     ______     ______   .___  ___. \n      |   _  \\   /  __  \\   /  __  \\  |   \\/   | \n      |  |_)  | |  |  |  | |  |  |  | |  \\  /  | \n      |   _  <  |  |  |  | |  |  |  | |  |\\/|  | \n      |  |_)  | |  `--'  | |  `--'  | |  |  |  | \n      |______/   \\______/   \\______/  |__|  |__| \n                                                 ")
	addr := config.RetServer().Address
	webs := config.RetServer().WebSocketPort
	data := [][]string{
		[]string{"webSocket", webs, "[ok]", "ws://127.0.0.1" + webs + "/ws"},
		[]string{"main", addr, "[ok]", "http://127.0.0.1" + addr},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Ports", "Status", "Request"})
	table.SetCenterSeparator("+")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor, tablewriter.FgGreenColor},
		tablewriter.Colors{},
	)
	table.AppendBulk(data)
	table.Render()
}
