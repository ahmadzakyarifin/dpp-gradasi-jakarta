// Kontrak API (docs-final) memakai `image_path` di response; halaman FE memakai `image_url`.
// Dipakai service layer agar semua halaman tetap membaca `image_url` tanpa menyentuh tiap halaman.
export function normalizeImage(item) {
  if (item && item.image_path && !item.image_url) {
    return { ...item, image_url: item.image_path }
  }
  return item
}
