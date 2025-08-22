// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"strings"
)

// node represents a single path segment in the routing tree.
// Each node can be either a static segment ("user") or a dynamic parameter (":id").
type node struct {
	segment   string           // path segment name
	param     bool             // true if this is a parameter node (":id")
	handlers  []HandlerFunc    // handlers executed if this route matches
	children  map[string]*node // child nodes for static segments
	paramNode *node            // child node dedicated to parameter segments
}

// addRoute registers a new route in the routing tree.
func (r *Engine) addRoute(method, path string, handlers ...HandlerFunc) {
	if r.trees[method] == nil {
		// init root node if not exists for this method
		r.trees[method] = &node{children: make(map[string]*node)}
	}
	current := r.trees[method]
	parts := splitClean(path)

	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			// check conflict: param cannot coexist with static child
			if _, exists := current.children[part]; exists {
				panic("conflict: param '" + part + "' collides with static route in " + method + " " + path)
			}

			// if no paramNode yet → create one
			if current.paramNode == nil {
				current.paramNode = &node{
					segment:  part[1:], // remove ":" to store only the param name
					param:    true,
					children: make(map[string]*node),
				}
			}

			// move deeper into paramNode
			current = current.paramNode
			continue
		} else {
			// check conflict: static cannot coexist with paramNode
			if current.paramNode != nil {
				panic("conflict: static '" + part + "' collides with param in " + method + " " + path)
			}

			// if child not exists → create one
			next := current.children[part]
			if next == nil {
				next = &node{segment: part, children: make(map[string]*node)}
				current.children[part] = next
			}

			// move deeper into static child
			current = next
		}
	}

	// after loop, current points to final node
	// check if handlers already exist → duplicate route
	if current.handlers != nil {
		panic("duplicate route detected: " + method + " " + path)
	}

	// assign handlers to this node
	current.handlers = handlers

	// add to route list for inspection/debug
	r.routeList = append(r.routeList, RouteInfo{
		Method: method,
		Path:   path,
	})
}

// findRoute searches for a matching route in the tree.
func (r *Engine) findRoute(method, path string) ([]HandlerFunc, map[string]string) {
	root := r.trees[method]
	if root == nil {
		// no route registered for this method
		return nil, nil
	}

	parts := splitClean(path)
	current := root
	var params map[string]string

	for _, part := range parts {
		// first try exact static match
		if next, ok := current.children[part]; ok {
			current = next
			continue
		}

		// fallback: check if paramNode exists
		if current.paramNode != nil {
			current = current.paramNode

			// allocate params map only when needed
			if params == nil {
				params = make(map[string]string)
			}

			// store actual value to param name
			params[current.segment] = part
			continue
		}

		// neither static nor param match → route not found
		return nil, nil
	}

	// reached final node, return handlers and params
	return current.handlers, params
}

func splitClean(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	raw := strings.Split(p, "/")
	out := raw[:0]
	for _, s := range raw {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
