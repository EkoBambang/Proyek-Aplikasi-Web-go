package main

//----------database----------------------
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "fmt"

//----------gorila------------------------
import "github.com/gorilla/mux"
import "github.com/gorilla/securecookie"

//----------http--------------------------
import "net/http"
import "html/template"
import "os"

//----------byncript----------------------
//import "golang.org/x/crypto/bcrypt"

//------------------------------connect to database---------------------------
var user_db, pass_db, login_db string

func connect() *sql.DB {
	var db, err = sql.Open("mysql", user_db+":"+pass_db+"@/"+login_db)
	err = db.Ping()
	if err != nil {
		fmt.Println("database tidak bisa dihubungi")
		os.Exit(0)

	}
	return db

}

//------------------------------login menu-------------------------------------
var router = mux.NewRouter()

var cookiehandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func halaman_login(res http.ResponseWriter, req *http.Request) {
	akses := namauser(req)
	if akses != "" {
		http.Redirect(res, req, "/penjualan", 301)
		return
	} else {

		http.ServeFile(res, req, "login.html")
	}
}
func mau_login(res http.ResponseWriter, req *http.Request) {
	username_login := req.FormValue("username_ku")
	password_login := req.FormValue("password_ku")
	jalur := "/"

	var nama, username, password string

	db := connect()
	defer db.Close()

	db.QueryRow("select nama,username,password from user where username=?", username_login).Scan(&nama, &username, &password)

	if password == password_login {
		setsession(username, res)
		jalur = "/penjualan"
	} else {

		jalur = "/"
	}
	http.Redirect(res, req, jalur, 302)
}

func halaman_index(res http.ResponseWriter, req *http.Request) {
	akses := namauser(req)
	if akses != "" {
		db := connect()
		defer db.Close()
		var ada_nama string
		db.QueryRow("select nama from user where username=?", akses).Scan(&ada_nama)

		halaman, _ := template.ParseFiles("utama.html")

		ident := map[string]string{
			"nama": ada_nama,
		}

		halaman.Execute(res, ident)

	} else {
		http.Redirect(res, req, "/", 301)
	}
}

func mau_logout(res http.ResponseWriter, req *http.Request) {
	clear_session(res)
	http.Redirect(res, req, "/", 301)
}

//----------------------------gagal handle------------------------------

func gagal_login(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "gagal_login.html")

}

//-----------------------------confirm-----------------------------------
func masih_login(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()
	var nama_yg_login string
	db.QueryRow("select nama from user where username=?", namauser(req)).Scan(&nama_yg_login)

	halaman, _ := template.ParseFiles("masih_login.html")

	ident := map[string]string{
		"nama": nama_yg_login,
	}

	halaman.Execute(res, ident)

}

//------------------------------------sesi--------------------------------------
func setsession(nama_user string, res http.ResponseWriter) {
	value := map[string]string{
		"name": nama_user,
	}
	if encoded, err := cookiehandler.Encode("session", value); err == nil {
		cookie_ku := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie_ku)
	}
}

func namauser(req *http.Request) (name_usernya string) {
	if cookie_ini, err := req.Cookie("session"); err == nil {
		nilai_cookie := make(map[string]string)
		if err = cookiehandler.Decode("session", cookie_ini.Value, &nilai_cookie); err == nil {
			name_usernya = nilai_cookie["name"]
		}

	}
	return name_usernya
}

func clear_session(res http.ResponseWriter) {
	bersihkan_cookie_ku := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, bersihkan_cookie_ku)
}

//----------------------------------------webb app------------------------------
func penjualan(res http.ResponseWriter, req *http.Request) {
	req_user := namauser(req)
	if req_user != "" {
		isi_nota(req)
		db := connect()
		defer db.Close()
		var nama_lengkap string
		db.QueryRow("select nama from user where username=?", req_user).Scan(&nama_lengkap)

		halaman, _ := template.ParseFiles("penjualan.html")

		data := map[string]string{
			"kd_barang":       "",
			"nama_barang":     "",
			"harga_barang":    "",
			"jumlah_barang":   "",
			"satuan_barang":   "",
			"diskon_barang":   "",
			"nama_user_login": nama_lengkap,
		}
		halaman.Execute(res, data)
	} else {
		http.Redirect(res, req, "/", 301)
	}
}
func mau_jual(res http.ResponseWriter, req *http.Request) {
	req_user := namauser(req)
	if req_user != "" {
		db := connect()
		defer db.Close()

		id_barang := req.FormValue("kd_barang")
		nama_barang := req.FormValue("nama_barang")
		harga := req.FormValue("harga_barang")
		quantity := req.FormValue("jumlah_barang")
		satuan := req.FormValue("satuan_barang")
		diskon := req.FormValue("diskon_barang")
		user_log := namauser(req)

		if nama_barang == "" {
			http.Redirect(res, req, "/data_kosong", 301)
			return
		}
		if harga == "" {
			http.Redirect(res, req, "/data_kosong", 301)
			return
		}

		_, err_1 := db.Exec("insert into nota (kd_barang,nama_barang,satuan,harga,jumlah,diskon,username) values (?,?,?,?,?,?/100,?)", id_barang, nama_barang, satuan, harga, quantity, diskon, user_log)
		if err_1 != nil {

		}
		_, err_2 := db.Exec("update barang set jumlah=jumlah-? where kd_barang=?", quantity, id_barang)
		if err_2 != nil {

		}
		http.Redirect(res, req, "/penjualan", 301)
	} else {
		http.Redirect(res, req, "/", 301)
	}
}
func cari_barang(res http.ResponseWriter, req *http.Request) {
	akses := namauser(req)
	if akses != "" {
		isi_nota(req)
		db := connect()
		defer db.Close()
		id_barang := req.FormValue("kd_barang")

		var kd, nama_barang, harga, jumlah, satuan, nama_user string
		db.QueryRow("select nama from user where username=?", akses).Scan(&nama_user)

		rows := db.QueryRow("select kd_barang,nama_barang,harga,jumlah,satuan from barang where kd_barang=?", id_barang)

		halaman, _ := template.ParseFiles("penjualan.html")
		rows.Scan(&kd, &nama_barang, &harga, &jumlah, &satuan)

		data := map[string]string{
			"kd_barang":       kd,
			"nama_barang":     nama_barang,
			"harga_barang":    harga,
			"jumlah_barang":   "",
			"satuan_barang":   satuan,
			"nama_user_login": nama_user,
		}
		halaman.Execute(res, data)
	} else {
		http.Redirect(res, req, "/", 301)
	}
}

func data_kosong(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "data_kosong.html")

}

func isi_nota(req *http.Request) {
	akses := namauser(req)

	var path = "data/tabel_nota.html"
	var path2 = "data/tabel_stok.html"

	os.Remove(path)
	os.Create(path)
	os.Remove(path2)
	os.Create(path2)

	var file, _ = os.OpenFile(path, os.O_RDWR, 0644)
	var file2, _ = os.OpenFile(path2, os.O_RDWR, 0644)

	defer file2.Close()
	defer file.Close()

	db := connect()
	defer db.Close()

	var no, kd_barang, nama_barang, satuan, harga, jumlah, diskon, total, semua_total string
	var kd_barang_stok, nama_barang_stok, satuan_stok, jumlah_stok, harga_stok, nama string

	rows_stok, _ := db.Query("select kd_barang,nama_barang,satuan,format(harga,'##,##0'),sum(jumlah) from barang group by kd_barang")
	rows, _ := db.Query("select @no := @no + 1 as no,kd_barang,nama_barang,satuan,format(harga,'##,##0'),jumlah,format(diskon*100,'##,##0'),format((harga*jumlah)-((harga*jumlah)*diskon),'##,##0') as total from nota join (select @no := 0) r where username=?", akses)
	db.QueryRow("select format(sum((harga-(harga*diskon))*jumlah),'##,##0') as total_beli from nota where username=?", akses).Scan(&semua_total)
	db.QueryRow("select nama from user where username=?", akses).Scan(&nama)

	file.WriteString("<table width='100%' border='0'>")
	file.WriteString("<tr><td>No</td><td>Nama Barang</td><td>Satuan</td><td>Harga</td><td>Jumlah</td><td>Diskon</td><td>Total</td></tr>")
	for rows.Next() {

		rows.Scan(&no, &kd_barang, &nama_barang, &satuan, &harga, &jumlah, &diskon, &total)
		file.WriteString("<tr><td>" + no + "</td><td>" + nama_barang + "</td><td>" + satuan + "</td><td> Rp " + harga + "</td><td>" + jumlah + " " + satuan + "</td><td>" + diskon + " % </td><td> Rp " + total + "</td></tr>")
	}
	file.WriteString("<tr><td colspan='7' height='30'></td></tr>")
	file.WriteString("<tr><td colspan='5'></td><td>Total semua : </td><td width='20%'> Rp <a id='total_semua_belanja'>" + semua_total + "</a></td></tr>")
	file.WriteString("<tr><td colspan='5'></td><td>Uang yg Diberikan : </td><td width='20%'> Rp <a id='uang_diberikan'>0</a></td></tr>")
	file.WriteString("<tr><td colspan='5'></td><td>kembalian : </td><td width='20%'> Rp <a id='kembalian'>0</a></td></tr>")
	file.WriteString("<tr><td colspan='5'></td><td>Nama user : </td><td> " + nama + "</td></tr>")
	file.WriteString("</table>")
	file.Sync()

	file2.WriteString("<table width='100% border='0'")
	file2.WriteString("<tr><td>Kode barang</td><td>Nama Barang</td><td>Satuan</td><td>Harga</td><td>Stok Tersedia</td>")
	for rows_stok.Next() {
		rows_stok.Scan(&kd_barang_stok, &nama_barang_stok, &satuan_stok, &harga_stok, &jumlah_stok)
		file2.WriteString("<tr><td>" + kd_barang_stok + "</td><td>" + nama_barang_stok + "</td><td>" + satuan_stok + "</td><td> Rp " + harga_stok + "</td><td>" + jumlah_stok + " " + satuan_stok + "</td>")

	}

	file2.WriteString("</table>")
	file2.Sync()
}

func hapus_nota(res http.ResponseWriter, req *http.Request) {
	nama_akses := namauser(req)
	if nama_akses != "" {
		db := connect()
		defer db.Close()

		db.Exec("delete from nota where username=?", nama_akses)

		http.Redirect(res, req, "/penjualan", 301)

	} else {
		http.Redirect(res, req, "/", 301)
	}

}

//----------------------------------------main func-----------------------------
func main() {

	fmt.Print("masukkan user database :")
	fmt.Scanln(&user_db)
	fmt.Print("masukkan pass database :")
	fmt.Scanln(&pass_db)
	fmt.Print("masukkan database database :")
	fmt.Scanln(&login_db)

	//----------------------------------------login logout----------------------

	router.HandleFunc("/", halaman_login)

	router.HandleFunc("/halaman_index", halaman_index)
	router.HandleFunc("/mau_login", mau_login)
	router.HandleFunc("/mau_logout", mau_logout)

	//---------------------------------------gagal and confirm func-------------

	router.HandleFunc("/gagal_login", gagal_login)
	router.HandleFunc("/masih_login", masih_login)

	//---------------------------------------web app----------------------------

	router.HandleFunc("/penjualan", penjualan)
	router.HandleFunc("/mau_jual", mau_jual)
	router.HandleFunc("/cari_barang", cari_barang)
	router.HandleFunc("/data_kosong", data_kosong)
	router.HandleFunc("/hapus_nota", hapus_nota)

	//-----------------------------------------server stuff--------------------------

	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))
	http.Handle("/", router)
	fmt.Println("running server via localhost:8080... ")
	http.ListenAndServe(":8080", nil)
}
