# Methods dan Helper dalam Go

Dalam proyek ini, kita menyediakan kumpulan materi dan contoh tentang method-method dan helper function dalam Go yang mirip dengan method-method di JavaScript seperti `map`, `filter`, `split`, `join`, dan lain-lain.

## Struktur Direktori

- `string_methods.md` - Method-method untuk manipulasi string
- `slice_methods.md` - Method-method untuk manipulasi slice (mirip dengan array methods di JavaScript)
- `map_methods.md` - Method-method untuk manipulasi map
- `number_methods.md` - Method-method untuk operasi matematika dan angka
- `struct_methods.md` - Method-method dengan receiver pada struct

## Konsep Utama

Dalam Go, tidak seperti JavaScript yang memiliki method bawaan pada tipe data seperti array dan string, kita perlu menggunakan:

1. **Package bawaan** seperti `strings`, `slices`, `math` untuk fungsi utilitas
2. **Fungsi generik** untuk membuat helper yang dapat bekerja dengan berbagai tipe data
3. **Method dengan receiver** untuk menambahkan perilaku ke struct
4. **Interface** untuk mencapai polimorfisme

## Perbedaan dengan JavaScript

| JavaScript | Go |
|------------|-----|
| `"hello".split(",")` | `strings.Split("hello", ",")` |
| `[1,2,3].map(x => x*2)` | Custom `Map` function |
| `[1,2,3].filter(x => x>1)` | Custom `Filter` function |
| `Math.max(1, 2)` | `math.Max(1, 2)` |
| Objek memiliki method bawaan | Menggunakan package atau implementasi sendiri |

## Contoh Penggunaan

Setiap file berisi contoh kode yang bisa langsung dijalankan dan diadaptasi sesuai kebutuhan. Contoh-contoh ini menunjukkan cara membuat helper function yang reusable dan idiomatik dalam Go.

## Lisensi

Materi ini dibuat untuk tujuan pembelajaran Go dan dapat digunakan secara bebas untuk pembelajaran.