const idn = require('idn-area-data');
async function test() {
  const prov = await idn.getProvinces();
  console.log(prov.slice(0, 2));
  const kab = await idn.getRegencies();
  console.log(kab.slice(0, 2));
}
test();
