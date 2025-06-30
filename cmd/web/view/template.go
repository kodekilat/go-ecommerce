package view

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

// cacheTemplates akan menyimpan template yang sudah di-parse di memori.
var cacheTemplates = make(map[string]*template.Template)

// Fungsi ini akan mendapatkan path absolut ke direktori root proyek.
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// b adalah path ke file ini (cmd/web/view/template.go)
	// kita naik 3 level untuk mendapatkan root proyek
	return filepath.Join(filepath.Dir(b), "../../..")
}

func Render(w http.ResponseWriter, tmplName string, data any) {
	// Untuk pengembangan, kita bisa nonaktifkan cache agar perubahan template langsung terlihat.
	// Untuk produksi, cache diaktifkan untuk performa.
	// Di sini kita buat sederhana dan selalu parse ulang.

	// Bangun path absolut ke file template
	rootPath := getProjectRoot()
	tmplPath := filepath.Join(rootPath, "web", "templates", tmplName)

	// Cek apakah file ada sebelum di-parse
	_, err := http.Dir(filepath.Dir(tmplPath)).Open(filepath.Base(tmplPath))
	if err != nil {
		log.Printf("Template file tidak ditemukan: %s", tmplPath)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse template
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmplPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Eksekusi template
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template %s: %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
