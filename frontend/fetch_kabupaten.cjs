const fs = require('fs');

async function fetchAll() {
  console.log("Fetching provinces...");
  const provRes = await fetch("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json");
  const provinces = await provRes.json();
  
  let allKabupatens = [];
  
  for (const prov of provinces) {
    console.log(`Fetching regencies for ${prov.name}...`);
    const regRes = await fetch(`https://www.emsifa.com/api-wilayah-indonesia/api/regencies/${prov.id}.json`);
    const regencies = await regRes.json();
    for (const reg of regencies) {
      // capitalize properly: "KABUPATEN BOGOR" -> "Kabupaten Bogor", "KOTA BEKASI" -> "Kota Bekasi"
      const name = reg.name.toLowerCase().replace(/\b\w/g, l => l.toUpperCase());
      allKabupatens.push(name);
    }
  }
  
  allKabupatens.sort();
  
  const content = `// Auto-generated list of all Regencies and Cities in Indonesia
export const INDONESIA_KABUPATEN = ${JSON.stringify(allKabupatens, null, 2)};
`;

  fs.writeFileSync('src/data/kabupaten.js', content);
  console.log(`Saved ${allKabupatens.length} kabupatens to src/data/kabupaten.js`);
}

fetchAll().catch(console.error);
