package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"time"
)

type Product struct {
	name string
	productid string
	desc string
	price float64
}

type User struct {
	admin bool
	username string
	name string
	address string
	dob time.Time
	cart []string
}

type App struct {
	loggedIn bool
	products	map[string]Product
	users []User
	creds map[string]string
	currentUser *User
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	a := App {
		loggedIn: false,
		products: make(map[string]Product, 0),
		users: make([]User, 0),
		creds: make(map[string]string, 0),
		currentUser: nil,
	}
	a.init()
	for {
		fmt.Println("login | signup | add_product | delete_product | list | add_to_cart | checkout | logout ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		a.parseAndExecute(text)
	}
}

func (a *App) init() {
	adminUser := User {
		admin: true,
		username: "admin",
		name: "admin",
		address: "nowhere",
		dob: time.Date(2022, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
	}
	a.creds["admin"] = "password1"
	a.users = append(a.users, adminUser)
}

func (a *App) parseAndExecute(s string) {
	switch s {
		case "login":
			a.login()
			break
		case "signup":
			a.signup()
			break
		case "add_product":
			a.addProduct()
			break
		case "delete_product":
			a.deleteProduct()
			break
		case "list":
			a.list()
			break
		case "add_to_cart":
			a.addToCart()
			break
		case "checkout":
			a.checkout()
			break
		case "logout":
			break
		default:
			fmt.Println("unrecognised command.")
	}	
}

func (a *App) login() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter username: ")
	un, _ := reader.ReadString('\n')
	un = strings.TrimSpace(un)
	fmt.Println("Enter password: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.TrimSpace(pw)
	if a.creds[un] == pw {
		a.loggedIn = true
		for i:=0; i<len(a.users); i++ {
			if a.users[i].username == un { a.currentUser = &a.users[i] }
		}
	} else {
		fmt.Println("Invalid username or password.")
	}
}

func (a *App) signup() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter username: ")
	un, _ := reader.ReadString('\n')
	un = strings.TrimSpace(un)
	if a.creds[un] != "" { fmt.Println("username exists."); return }
	fmt.Println("Enter password: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.TrimSpace(pw)
	fmt.Println("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.TrimSpace(addr)
	fmt.Println("Enter dob(dd-mm-yyyy): ")
	dobString, _ := reader.ReadString('\n')
	dobString = strings.TrimSpace(dobString)
	dateStrings := strings.Split(dobString, "-")
	year, yErr := strconv.Atoi(dateStrings[0])
	mon, mErr := strconv.Atoi(dateStrings[1])
	day, dErr := strconv.Atoi(dateStrings[2])
	var dob time.Time
	if (yErr==nil && mErr==nil && dErr==nil) {
		dob = time.Date(year, time.Month(mon), day, 0, 0, 0, 0, time.UTC)
	} else {
		fmt.Println("can't process date.")
		return
	}
	user := User {
		username: un,
		name: name,
		address: addr,
		dob: dob,
	}
	a.users = append(a.users, user)
	a.creds[un] = pw
}

func (a *App) addProduct() {
	if !a.loggedIn { fmt.Println("not logged in."); return }
	if !a.currentUser.admin { fmt.Println("not an admin"); return }
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter product id: ")
	pid, _ := reader.ReadString('\n')
	pid = strings.TrimSpace(pid)
	fmt.Println("Enter product name: ")
	pname, _ := reader.ReadString('\n')
	pname = strings.TrimSpace(pname)
	fmt.Println("Enter description: ")
	pdesc, _ := reader.ReadString('\n')
	pdesc = strings.TrimSpace(pdesc)
	fmt.Println("Enter price: ")
	ppricestr, _ := reader.ReadString('\n')
	ppricestr = strings.TrimSpace(ppricestr)
	price, err := strconv.ParseFloat(ppricestr, 32)
	if err != nil { 
		fmt.Println("can't parse price")
		return
	}
	p := Product {
		productid: pid,
		name: pname,
		desc: pdesc,
		price: price,
	}
	a.products[pid] = p
}

func (a *App) deleteProduct() {
	if !a.loggedIn { fmt.Println("not logged in."); return }
	if !a.currentUser.admin { fmt.Println("not an admin"); return }
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter product id: ")
	pid, _ := reader.ReadString('\n')
	pid = strings.TrimSpace(pid)
	delete(a.products, pid)
			a.list()
}

func (a *App) list() {
	if !a.loggedIn { fmt.Println("not logged in."); return }
	for pid, p := range a.products {
		fmt.Println(pid + " " + p.name + " " + p.desc + " " + strconv.FormatFloat(p.price, 'E', -1, 64))
	}
}

func (a *App) addToCart() {
	if !a.loggedIn { fmt.Println("not logged in."); return }
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter product id: ")
	pid, _ := reader.ReadString('\n')
	pid = strings.TrimSpace(pid)
	a.currentUser.cart = append(a.currentUser.cart, pid)
}

func (a *App) checkout() {
	if !a.loggedIn { fmt.Println("not logged in."); return }
	var total float64
	for i:=0; i<len(a.currentUser.cart); i++ {
		p := a.products[a.currentUser.cart[i]]
		fmt.Println(p.productid + " " + p.name + " " + strconv.FormatFloat(p.price, 'E', -1, 64))
		total += p.price
	}
	fmt.Println("total: " + strconv.FormatFloat(total, 'E', -1, 64))
} 

func (a *App) logout() {
	a.loggedIn = false
}
