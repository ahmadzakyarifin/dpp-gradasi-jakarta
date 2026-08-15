const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ headless: "new", args: ['--no-sandbox'] });
  const page = await browser.newPage();
  page.on('console', msg => console.log('PAGE LOG:', msg.text()));
  page.on('pageerror', error => console.log('PAGE ERROR:', error.message));
  
  await page.goto('http://localhost:5173/kepengurusan?tab=Pengurus+Provinsi', { waitUntil: 'networkidle0', timeout: 10000 }).catch(e => console.log('Nav error:', e.message));
  
  await browser.close();
})();
