import os
import re

html_files = [f for f in os.listdir('admin') if f.endswith('.html')]

target = r'<div class="flex items-center gap-4">\s*<div class="flex items-center gap-2">\s*<img src="https://ui-avatars\.com/api/\?name=Super\+Admin&background=0D8ABC&color=fff" class="w-8 h-8 rounded-full" alt="Admin">\s*<span class="text-sm font-medium text-gray-700">Super Admin</span>\s*</div>\s*</div>'

replacement = """<div class="flex items-center gap-4">
                <div class="relative" x-data="{ userMenuOpen: false }" @click.away="userMenuOpen = false">
                    <button @click="userMenuOpen = !userMenuOpen" class="flex items-center gap-2 hover:bg-gray-100 p-2 rounded-lg transition">
                        <img src="https://ui-avatars.com/api/?name=Admin&background=0D8ABC&color=fff" class="w-8 h-8 rounded-full" alt="Admin" id="headerProfileImg">
                        <span class="text-sm font-medium text-gray-700" id="headerProfileName">Admin</span>
                        <i class="ph ph-caret-down text-gray-500"></i>
                    </button>
                    <div x-show="userMenuOpen" x-transition 
                         class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 py-1 z-50 hidden"
                         :class="{'hidden': !userMenuOpen}">
                        <a href="profile.html" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 hover:text-brand-600">
                            <i class="ph ph-user-circle mr-2"></i> Profil Saya
                        </a>
                        <div class="border-t border-gray-100 my-1"></div>
                        <a href="../login.html" onclick="localStorage.removeItem('access_token')" class="block px-4 py-2 text-sm text-red-600 hover:bg-red-50">
                            <i class="ph ph-sign-out mr-2"></i> Logout
                        </a>
                    </div>
                </div>
            </div>"""

for file in html_files:
    filepath = os.path.join('admin', file)
    with open(filepath, 'r') as f:
        content = f.read()
    
    new_content = re.sub(target, replacement, content)
    
    if new_content != content:
        with open(filepath, 'w') as f:
            f.write(new_content)
        print(f"Updated {file}")
