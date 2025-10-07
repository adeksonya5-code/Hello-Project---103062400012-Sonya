package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const usersFile = "users.txt"
const booksFile = "books.txt"
//ini program awal tugas besar semester 2

type User struct {
	IDPengguna   int
	NamaPengguna string
	NoHP         string
}

type Buku struct {
	IDBuku    int
	JudulBuku string
	Pengarang string
	HargaBuku float64
	Stok      int
}

type Pembelian struct {
	IDPembelian int
	IDUser      int
	IDBuku      int
	Jumlah      int
}

var daftarBuku []Buku
var nextIDBuku = 1
var riwayatPembelian []Pembelian
var nextIDPembelian = 1
var nextIDPengguna = 1
var shoppingCarts = make(map[int]map[int]int)

//(fungsi readFileLines dan writeFileLines tidak berubah) ...
func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeFileLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
// (fungsi loadBooksFromFile dan saveBooksToFile tidak berubah) ...
func loadBooksFromFile() {
	lines, err := readFileLines(booksFile)
	if err != nil {
		fmt.Println("Error loading books:", err)
		return
	}

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 5 {
			id, _ := strconv.Atoi(parts[0])
			harga, _ := strconv.ParseFloat(parts[3], 64)
			stok, _ := strconv.Atoi(parts[4])
			buku := Buku{
				IDBuku:    id,
				JudulBuku: parts[1],
				Pengarang: parts[2],
				HargaBuku: harga,
				Stok:      stok,
			}
			daftarBuku = append(daftarBuku, buku)
			if id >= nextIDBuku {
				nextIDBuku = id + 1
			}
		}
	}
}

func saveBooksToFile() {
	var lines []string
	for _, buku := range daftarBuku {
		line := fmt.Sprintf("%d,%s,%s,%.2f,%d", buku.IDBuku, buku.JudulBuku, buku.Pengarang, buku.HargaBuku, buku.Stok)
		lines = append(lines, line)
	}
	_ = writeFileLines(booksFile, lines)
}

// ===== User Handling =====
// ... (fungsi loadUsersFromFile dan saveUsersToFile tidak berubah) ...
func loadUsersFromFile() ([]User, error) {
	var users []User
	lines, err := readFileLines(usersFile)
	if err != nil {
		return nil, fmt.Errorf("error reading users file: %w", err)
	}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 {
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println("Warning: invalid user ID:", line)
				continue
			}
			user := User{IDPengguna: id, NamaPengguna: parts[1], NoHP: parts[2]}
			users = append(users, user)
			if id >= nextIDPengguna {
				nextIDPengguna = id + 1
			}
		} else if len(line) > 0 {
			fmt.Println("Warning: invalid user format:", line)
		}
	}
	return users, nil
}

func saveUsersToFile(users []User) error {
	var lines []string
	for _, user := range users {
		lines = append(lines, strings.Join([]string{
			strconv.Itoa(user.IDPengguna),
			user.NamaPengguna,
			user.NoHP,
		}, ","))
	}
	return writeFileLines(usersFile, lines)
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func adminAddUserInteractive() {
	fmt.Println("\n--- Admin: Add New User ---")
	namaPengguna := getInput("Nama Pengguna: ")
	noHP := getInput("Nomor HP: ")

	users, err := loadUsersFromFile()
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}
	newUser := User{IDPengguna: nextIDPengguna, NamaPengguna: namaPengguna, NoHP: noHP}
	users = append(users, newUser)
	if err := saveUsersToFile(users); err != nil {
		fmt.Println("Error saving users:", err)
		return
	}
	fmt.Printf("Admin: User '%s' (ID: %d) created and saved.\n", namaPengguna, newUser.IDPengguna)
	nextIDPengguna++
}

func adminListUsers() {
	fmt.Println("\n--- Admin: User List ---")
	users, err := loadUsersFromFile()
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}
	if len(users) == 0 {
		fmt.Println("Admin: No users in file.")
		return
	}
	for _, user := range users {
		fmt.Printf("ID: %d, Nama: %s, No. HP: %s\n", user.IDPengguna, user.NamaPengguna, user.NoHP)
	}
	fmt.Println("------------------------")
}

func sequentialSearchUser(namaCari string) {
	fmt.Printf("\n--- Admin: Sequential Search User (Nama: '%s') ---\n", namaCari)
	lines, err := readFileLines(usersFile)
	if err != nil {
		fmt.Println("Error reading users file:", err)
		return
	}

	found := false
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 {
			nama := parts[1]
			if strings.Contains(strings.ToLower(nama), strings.ToLower(namaCari)) {
				fmt.Println("Ditemukan:", line)
				found = true
			}
		}
	}

	if !found {
		fmt.Printf("Tidak ada pengguna dengan nama mengandung '%s' ditemukan.\n", namaCari)
	}
	fmt.Println("----------------------------------------------------")
}

func adminEditUserInteractive() {
	fmt.Println("\n--- Admin: Edit/Delete User ---")
	idEditStr := getInput("Masukkan ID Pengguna yang ingin diedit atau dihapus: ")
	idEdit, err := strconv.Atoi(idEditStr)
	if err != nil {
		fmt.Println("Invalid ID format.")
		return
	}

	users, err := loadUsersFromFile()
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}

	foundIndex := -1
	for i, user := range users {
		if user.IDPengguna == idEdit {
			foundIndex = i
			fmt.Printf("\n--- User ID: %d ---\n", idEdit)
			fmt.Printf("Nama saat ini: %s\n", user.NamaPengguna)
			fmt.Printf("No. HP saat ini: %s\n", user.NoHP)
			break
		}
	}

	if foundIndex == -1 {
		fmt.Printf("Pengguna dengan ID %d tidak ditemukan.\n", idEdit)
		return
	}

	actionChoice := getInput("\nApa yang ingin Anda lakukan? (1. Edit, 2. Hapus): ")
	switch actionChoice {
	case "1":
		namaBaru := getInput("Masukkan Nama Pengguna baru (kosongkan untuk tidak mengubah): ")
		if namaBaru != "" {
			users[foundIndex].NamaPengguna = namaBaru
		}
		noHPBaru := getInput("Masukkan Nomor HP baru (kosongkan untuk tidak mengubah): ")
		if noHPBaru != "" {
			users[foundIndex].NoHP = noHPBaru
		}
		if err := saveUsersToFile(users); err != nil {
			fmt.Println("Error saving users:", err)
		} else {
			fmt.Printf("\nUser ID %d berhasil diupdate.\n", idEdit)
		}
	case "2":
		users = append(users[:foundIndex], users[foundIndex+1:]...)
		if err := saveUsersToFile(users); err != nil {
			fmt.Println("Error saving users:", err)
		} else {
			fmt.Printf("\nUser ID %d berhasil dihapus.\n", idEdit)
		}
	default:
		fmt.Println("Pilihan tidak valid.")
	}
}

func adminTambahBuku() {
	fmt.Println("\n--- Admin: Tambah Buku ---")
	judul := getInput("Judul Buku: ")
	pengarang := getInput("Pengarang: ")
	hargaStr := getInput("Harga Buku: ")
	stokStr := getInput("Stok Buku: ")

	harga, _ := strconv.ParseFloat(hargaStr, 64)
	stok, _ := strconv.Atoi(stokStr)

	buku := Buku{
		IDBuku:    nextIDBuku,
		JudulBuku: judul,
		Pengarang: pengarang,
		HargaBuku: harga,
		Stok:      stok,
	}
	daftarBuku = append(daftarBuku, buku)
	nextIDBuku++
	saveBooksToFile()
	fmt.Println("Buku berhasil ditambahkan.")
}

func adminLihatDaftarBuku() {
	fmt.Println("\n--- Daftar Buku ---")
	if len(daftarBuku) == 0 {
		fmt.Println("Belum ada buku.")
		return
	}
	for _, buku := range daftarBuku {
		fmt.Printf("ID: %d | Judul: %s | Pengarang: %s | Harga: %.2f | Stok: %d\n",
			buku.IDBuku, buku.JudulBuku, buku.Pengarang, buku.HargaBuku, buku.Stok)
	}
}

func adminEditStokBuku() {
	fmt.Println("\n--- Admin: Edit Stok Buku ---")
	idStr := getInput("Masukkan ID Buku: ")
	id, _ := strconv.Atoi(idStr)

	for i, buku := range daftarBuku {
		if buku.IDBuku == id {
			fmt.Printf("Stok saat ini: %d\n", buku.Stok)
			stokBaruStr := getInput("Masukkan stok baru: ")
			stokBaru, _ := strconv.Atoi(stokBaruStr)
			daftarBuku[i].Stok = stokBaru
			saveBooksToFile()
			fmt.Println("Stok buku diperbarui.")
			return
		}
	}
	fmt.Println("Buku tidak ditemukan.")
}

func hapusBuku() {
	fmt.Println("\n--- Admin: Hapus Buku ---")
	idStr := getInput("Masukkan ID buku yang ingin dihapus: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("ID tidak valid.")
		return
	}

	index := -1
	for i, buku := range daftarBuku {
		if buku.IDBuku == id {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Buku tidak ditemukan.")
		return
	}

	// Konfirmasi sebelum menghapus
	fmt.Printf("Apakah Anda yakin ingin menghapus buku '%s'? (y/n): ", daftarBuku[index].JudulBuku)
	konfirmasi := getInput("")
	if strings.ToLower(konfirmasi) != "y" {
		fmt.Println("Penghapusan dibatalkan.")
		return
	}

	// Hapus buku dari slice
	daftarBuku = append(daftarBuku[:index], daftarBuku[index+1:]...)
	saveBooksToFile()

	fmt.Println("Buku berhasil dihapus.")
}


func adminLihatRiwayatPembelian() {
	fmt.Println("\n--- Riwayat Pembelian ---")
	for _, p := range riwayatPembelian {
		fmt.Printf("ID Pembelian: %d | User ID: %d | Buku ID: %d | Jumlah: %d\n",
			p.IDPembelian, p.IDUser, p.IDBuku, p.Jumlah)
	}
}
// ===== User Functions =====
func userTambahKeKeranjang(user User) {
	fmt.Println("\n--- User: Tambah ke Keranjang ---")
	adminLihatDaftarBuku()
	idBukuStr := getInput("Masukkan ID Buku yang ingin ditambahkan ke keranjang: ")
	idBuku, err := strconv.Atoi(idBukuStr)
	if err != nil {
		fmt.Println("ID Buku tidak valid.")
		return
	}

	jumlahStr := getInput("Masukkan jumlah: ")
	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil {
		fmt.Println("Jumlah tidak valid.")
		return
	}

	bukuDitemukan := false
	for _, buku := range daftarBuku {
		if buku.IDBuku == idBuku {
			bukuDitemukan = true
			if buku.Stok < jumlah {
				fmt.Printf("Stok tidak mencukupi untuk buku '%s'. Stok tersedia: %d\n", buku.JudulBuku, buku.Stok)
				return
			}
			break
		}
	}
	if !bukuDitemukan {
		fmt.Println("Buku dengan ID tersebut tidak ditemukan.")
		return
	}

	userID := user.IDPengguna
	if _, exists := shoppingCarts[userID]; !exists {
		shoppingCarts[userID] = make(map[int]int)
	}
	shoppingCarts[userID][idBuku] += jumlah
	fmt.Printf("%d eksemplar buku dengan ID %d berhasil ditambahkan ke keranjang.\n", jumlah, idBuku)
}

func userLihatKeranjang(user User) {
	fmt.Println("\n--- User: Keranjang Belanja Anda ---")
	cart, exists := shoppingCarts[user.IDPengguna]
	if !exists || len(cart) == 0 {
		fmt.Println("Keranjang belanja Anda kosong.")
		return
	}

	var totalHargaKeranjang float64
	for bookID, quantity := range cart {
		for _, buku := range daftarBuku {
			if buku.IDBuku == bookID {
				hargaTotalItem := buku.HargaBuku * float64(quantity)
				fmt.Printf("ID: %d | Judul: %s | Jumlah: %d | Harga per unit: %.2f | Total: %.2f\n",
					buku.IDBuku, buku.JudulBuku, quantity, buku.HargaBuku, hargaTotalItem)
				totalHargaKeranjang += hargaTotalItem
				break
			}
		}
	}
	fmt.Printf("-----------------------------------\n")
	fmt.Printf("Total Harga Keranjang: %.2f\n", totalHargaKeranjang)
}

func userCheckout(user User) {
	fmt.Println("\n--- User: Checkout ---")
	cart, exists := shoppingCarts[user.IDPengguna]
	if !exists || len(cart) == 0 {
		fmt.Println("Keranjang belanja Anda kosong. Tidak ada yang di-checkout.")
		return
	}

	totalHargaCheckout := 0.0
	pembelianBerhasil := true
	var pembelianItems []Pembelian

	fmt.Println("\n--- Detail Checkout ---")
	for bookID, quantity := range cart {
		bukuIndex := -1
		var bukuYangDibeli *Buku
		for i, buku := range daftarBuku {
			if buku.IDBuku == bookID {
				bukuIndex = i
				bukuYangDibeli = &daftarBuku[i]
				break
			}
		}

		if bukuIndex == -1 {
			fmt.Printf("Peringatan: Buku dengan ID %d tidak ditemukan.\n", bookID)
			pembelianBerhasil = false
			continue
		}

		if bukuYangDibeli.Stok < quantity {
			fmt.Printf("Maaf, stok untuk buku '%s' tidak mencukupi (tersedia: %d, ingin dibeli: %d).\n",
				bukuYangDibeli.JudulBuku, bukuYangDibeli.Stok, quantity)
			pembelianBerhasil = false
			continue
		}

		hargaTotalItem := bukuYangDibeli.HargaBuku * float64(quantity)
		fmt.Printf("Judul: %s | Jumlah: %d | Harga per unit: %.2f | Total: %.2f\n",
			bukuYangDibeli.JudulBuku, quantity, bukuYangDibeli.HargaBuku, hargaTotalItem)
		totalHargaCheckout += hargaTotalItem

		// Siapkan data pembelian (belum dicatat permanen)
		pembelianItems = append(pembelianItems, Pembelian{
			IDPembelian: 0, // Akan diisi nanti
			IDUser:      user.IDPengguna,
			IDBuku:      bookID,
			Jumlah:      quantity,
		})
	}

	fmt.Printf("---------------------------\n")
	fmt.Printf("Total Harga Checkout: %.2f\n", totalHargaCheckout)

	if pembelianBerhasil && len(pembelianItems) > 0 {
		konfirmasi := getInput("Lanjutkan ke pembayaran? (y/n): ")
		if strings.ToLower(konfirmasi) == "y" {
			for _, item := range pembelianItems {
				// Kurangi stok
				for i := range daftarBuku {
					if daftarBuku[i].IDBuku == item.IDBuku {
						daftarBuku[i].Stok -= item.Jumlah
						break
					}
				}
				// Catat pembelian
				item.IDPembelian = nextIDPembelian
				riwayatPembelian = append(riwayatPembelian, item)
				nextIDPembelian++
			}
			saveBooksToFile()
			fmt.Println("Checkout berhasil! Terima kasih atas pembelian Anda.")
			// Kosongkan keranjang belanja setelah berhasil checkout
			delete(shoppingCarts, user.IDPengguna)
		} else {
			fmt.Println("Checkout dibatalkan.")
		}
	} else {
		fmt.Println("Tidak dapat melanjutkan checkout karena ada masalah dengan item di keranjang.")
	}
}

// ===== User Menu =====

func userMenu(user User) {
	for {
		fmt.Printf("\n--- User Menu (Halo %s) ---\n", user.NamaPengguna)
		fmt.Println("1. Tambah ke Keranjang")
		fmt.Println("2. Lihat Daftar Buku")
		fmt.Println("3. Lihat Keranjang")
		fmt.Println("4. Checkout")
		fmt.Println("5. Lihat Riwayat Pembelian")
		fmt.Println("6. Exit")
		choice := getInput("Enter your choice: ")

		switch choice {
		case "1":
			userTambahKeKeranjang(user)
		case "2":
			adminLihatDaftarBuku()
		case "3":
			userLihatKeranjang(user)
		case "4":
			userCheckout(user)
		case "5":
			fmt.Println("\n--- Riwayat Pembelian Anda ---")
			found := false
			for _, p := range riwayatPembelian {
				if p.IDUser == user.IDPengguna {
					fmt.Printf("ID Pembelian: %d | Buku ID: %d | Jumlah: %d\n",
						p.IDPembelian, p.IDBuku, p.Jumlah)
					found = true
				}
			}
			if !found {
				fmt.Println("Belum ada pembelian.")
			}
		case "6":
			fmt.Println("Keluar dari user menu.")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// ===== Admin Menu =====
// ... (tidak ada perubahan signifikan) ...

func adminMainMenu() {
	var choice string
	for {
		fmt.Println("\n--- Admin Menu ---")
		fmt.Println("1. Menambahkan  Pengguna")
		fmt.Println("2. List Pengguna")
		fmt.Println("3. Edit/Hapus Pengguna")
		fmt.Println("4. Mencari Pengguna berdasarkan Nama")
		fmt.Println("5. Menambahkan Buku")
		fmt.Println("6. List Buku")
		fmt.Println("7. Edit Stok Buku")
		fmt.Println("8. Riwayat Pembelian")
		fmt.Println("9. Hapus Buku")
		fmt.Println("10. Exit")
		fmt.Print("Masukan Pilihan Kamu: ")

		_, _ = fmt.Scanln(&choice)
		switch choice {
		case "1":
			adminAddUserInteractive()
		case "2":
			adminListUsers()
		case "3":
			adminEditUserInteractive()
		case "4":
			namaCari := getInput("Masukkan nama pengguna yang ingin dicari: ")
			sequentialSearchUser(namaCari)
		case "5":
			adminTambahBuku()
		case "6":
			adminLihatDaftarBuku()
		case "7":
			adminEditStokBuku()
		case "8":
			adminLihatRiwayatPembelian()
		case "9":
			hapusBuku()
		case "10":
			fmt.Println("Exiting Admin Panel. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

// ===== Main Entry =====
// (tidak ada perubahan signifikan) 

func main() {
	users, err := loadUsersFromFile()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading users:", err)
	}
	loadBooksFromFile()
	if len(users) > 0 {
		lastID := 0
		for _, u := range users {
			if u.IDPengguna > lastID {
				lastID = u.IDPengguna
			}
		}
		nextIDPengguna = lastID + 1
	} else {
		nextIDPengguna = 1
	}

	for {
		fmt.Println("--- Selamat Datang di Toko Buku ---")
	inputKode := getInput("Masuk sebagai : ")

	switch inputKode {
	case "Admin":
		password := getInput("Masukkan password admin: ")
			if password == "1" {
				fmt.Println("\n--- Panel Admin ---")
				adminMainMenu()
			} else {
				fmt.Println("Password salah. Akses ditolak.")
			}
	case "Pengguna":
		fmt.Println("\n--- Masuk sebagai Pengguna ---")
		users, _ := loadUsersFromFile()
		if len(users) == 0 {
			fmt.Println("Belum ada pengguna yang terdaftar.")
		} else {
			idStr := getInput("Masukkan ID Pengguna Anda: ")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("ID tidak valid.")
			} else {
				var foundUser *User
				for _, user := range users {
					if user.IDPengguna == id {
						foundUser = &user
						break
					}
				}
				if foundUser == nil {
					fmt.Println("Pengguna tidak ditemukan.")
				} else {
					userMenu(*foundUser)
				}
			}
		}
	default:
		fmt.Println("Pilihan tidak valid.")
	}
	}
}