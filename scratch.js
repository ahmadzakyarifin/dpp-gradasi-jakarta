const fs = require('fs');
const files = fs.readdirSync('admin').filter(f => f.endsWith('.html'));

const paginationUI = `
                    <!-- Pagination -->
                    <div x-show="filteredItems.length > 0" class="bg-gray-50 border-t border-gray-200 px-4 py-3 flex items-center justify-between sm:px-6 rounded-b-xl">
                        <div class="text-sm text-gray-500">
                            Menampilkan <span class="font-medium text-gray-900" x-text="(currentPage - 1) * pageSize + 1"></span> sampai <span class="font-medium text-gray-900" x-text="Math.min(currentPage * pageSize, filteredItems.length)"></span> dari <span class="font-medium text-gray-900" x-text="filteredItems.length"></span> hasil
                        </div>
                        <div class="flex gap-1">
                            <button @click="prevPage" :disabled="currentPage === 1" class="px-3 py-1 border border-gray-200 bg-white text-gray-500 rounded hover:bg-gray-50 disabled:opacity-50">Prev</button>
                            
                            <template x-for="p in totalPages" :key="p">
                                <button @click="goToPage(p)" class="px-3 py-1 border rounded" :class="currentPage === p ? 'border-brand-500 bg-brand-50 text-brand-600 font-medium' : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50'" x-text="p"></button>
                            </template>

                            <button @click="nextPage" :disabled="currentPage === totalPages" class="px-3 py-1 border border-gray-200 bg-white text-gray-500 rounded hover:bg-gray-50 disabled:opacity-50">Next</button>
                        </div>
                    </div>`;

const alpineLogic = `
                currentPage: 1,
                pageSize: 5,

                get totalPages() {
                    return Math.ceil(this.filteredItems.length / this.pageSize) || 1;
                },

                get paginatedItems() {
                    const start = (this.currentPage - 1) * this.pageSize;
                    return this.filteredItems.slice(start, start + this.pageSize);
                },

                nextPage() {
                    if (this.currentPage < this.totalPages) this.currentPage++;
                },

                prevPage() {
                    if (this.currentPage > 1) this.currentPage--;
                },
                
                goToPage(p) {
                    this.currentPage = p;
                },
`;

for (const file of files) {
  if (file === 'index.html' || file === 'login.html') continue;
  let content = fs.readFileSync('admin/' + file, 'utf8');

  // 1. Remove existing static pagination if any
  content = content.replace(/<!-- Pagination \(Static Mock\) -->[\s\S]*?<\/div>\s*<\/div>\s*<\/div>/g, '</div>\n</div>');
  content = content.replace(/<!-- Pagination -->[\s\S]*?<\/div>\s*<\/div>\s*<\/div>/g, '</div>\n</div>'); // for activity-log.html which has existing 

  // 2. Change x-for
  if (file === 'activity-log.html') {
    content = content.replace(/x-for="log in filteredLogs"/g, 'x-for="log in paginatedItems"');
    // For activity log, filteredItems is called filteredLogs. I will normalize it or adjust alpineLogic for it.
    // Actually, in activity-log, let's just make it filteredItems to unify, or change alpineLogic to filteredLogs
    content = content.replace(/filteredLogs/g, 'filteredItems');
    content = content.replace(/logs:/g, 'allData:');
  } else {
    content = content.replace(/x-for="item in filteredItems"/g, 'x-for="item in paginatedItems"');
  }

  // 3. Inject pagination UI after Empty State div
  // The empty state usually ends with </div> \n </div>
  // We'll just look for:
  content = content.replace(/(<!-- Empty State -->[\s\S]*?<\/div>\s*<\/div>)/, `$1\n${paginationUI}`);

  // 4. Inject logic into Alpine
  if (!content.includes('currentPage:')) {
    // find 'currentTab:' or 'filterStatus:'
    content = content.replace(/currentTab:[^,]*,/, (m) => m + alpineLogic);
  }

  // 5. Reset page when tab changes
  content = content.replace(/this\.\$watch\('currentTab', value => \{/, "this.$watch('currentTab', value => {\n                        this.currentPage = 1;");
  content = content.replace(/this\.\$watch\('filterStatus', value => \{/, "this.$watch('filterStatus', value => {\n                        this.currentPage = 1;");

  fs.writeFileSync('admin/' + file, content);
  console.log('Processed', file);
}
