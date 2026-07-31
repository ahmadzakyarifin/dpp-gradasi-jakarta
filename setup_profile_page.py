import re

with open('admin/profile.html', 'r') as f:
    content = f.read()

# Replace title
content = re.sub(r'<title>Dashboard - Admin DPP GRADASI</title>', '<title>Profil Saya - Admin DPP GRADASI</title>', content)
# Replace active menu
content = re.sub(r'href="index.html" class="flex items-center px-3 py-2.5 rounded-lg bg-brand-600 text-white shadow-sm transition-colors"', 
                 'href="index.html" class="flex items-center px-3 py-2.5 rounded-lg text-white/70 hover:bg-white/10 hover:text-white transition-colors"', content)
# Replace header title
content = re.sub(r'<h1 class="font-heading font-semibold text-gray-800 text-lg">Dashboard</h1>', 
                 '<h1 class="font-heading font-semibold text-gray-800 text-lg">Profil Saya</h1>', content)

# Replace content
profile_content = """<div class="flex-1 overflow-auto p-8">
            <div class="max-w-4xl mx-auto space-y-6" x-data="profileData()">
                
                <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
                    <div class="p-6 border-b border-gray-200">
                        <h2 class="text-lg font-semibold text-gray-800">Informasi Pribadi</h2>
                        <p class="text-sm text-gray-500">Perbarui nama, email, dan foto profil Anda.</p>
                    </div>
                    <div class="p-6">
                        <form @submit.prevent="updateProfile">
                            <div class="flex items-center gap-6 mb-6">
                                <img :src="formData.photo_preview || 'https://ui-avatars.com/api/?name=Admin&background=0D8ABC&color=fff'" class="w-20 h-20 rounded-full object-cover border border-gray-200" alt="Foto">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-2">Foto Profil</label>
                                    <input type="file" @change="handleFileChange" accept="image/*" class="text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-brand-50 file:text-brand-700 hover:file:bg-brand-100">
                                </div>
                            </div>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
                                    <input type="text" x-model="formData.name" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm">
                                </div>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">Alamat Email</label>
                                    <input type="email" x-model="formData.email" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm">
                                </div>
                            </div>
                            <div class="flex justify-end">
                                <button type="submit" class="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition shadow-sm text-sm font-medium" :disabled="loading">
                                    <span x-show="!loading">Simpan Perubahan</span>
                                    <span x-show="loading">Menyimpan...</span>
                                </button>
                            </div>
                        </form>
                    </div>
                </div>

                <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
                    <div class="p-6 border-b border-gray-200">
                        <h2 class="text-lg font-semibold text-gray-800">Ubah Password</h2>
                        <p class="text-sm text-gray-500">Pastikan akun Anda aman dengan password yang kuat.</p>
                    </div>
                    <div class="p-6">
                        <form @submit.prevent="updatePassword">
                            <div class="space-y-4 mb-6">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">Password Lama</label>
                                    <input type="password" x-model="passwordData.old_password" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm">
                                </div>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">Password Baru</label>
                                    <input type="password" x-model="passwordData.new_password" required minlength="6" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm">
                                </div>
                            </div>
                            <div class="flex justify-end">
                                <button type="submit" class="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition shadow-sm text-sm font-medium" :disabled="loadingPwd">
                                    <span x-show="!loadingPwd">Ubah Password</span>
                                    <span x-show="loadingPwd">Memproses...</span>
                                </button>
                            </div>
                        </form>
                    </div>
                </div>

            </div>
        </div>"""

content = re.sub(r'<div class="flex-1 overflow-auto p-8">.*</div>\s*</main>', profile_content + '\n    </main>', content, flags=re.DOTALL)

# Replace Alpine data script
script_content = """<script>
        document.addEventListener('alpine:init', () => {
            Alpine.data('profileData', () => ({
                formData: { name: '', email: '', photo: null, photo_preview: '' },
                passwordData: { old_password: '', new_password: '' },
                loading: false,
                loadingPwd: false,
                
                init() {
                    this.fetchProfile();
                },
                
                fetchProfile() {
                    const token = localStorage.getItem('access_token');
                    if (!token) return window.location.href = '../login.html';
                    
                    fetch('/api/v1/profile', {
                        headers: { 'Authorization': 'Bearer ' + token, 'Accept': 'application/json' }
                    })
                    .then(r => r.json())
                    .then(res => {
                        if (res.success && res.data) {
                            this.formData.name = res.data.name;
                            this.formData.email = res.data.email;
                            this.formData.photo_preview = res.data.photo_path ? ('http://127.0.0.1:8080' + res.data.photo_path) : '';
                            
                            // Update header info as well
                            document.getElementById('headerProfileName').innerText = res.data.name;
                            if (res.data.photo_path) {
                                document.getElementById('headerProfileImg').src = 'http://127.0.0.1:8080' + res.data.photo_path;
                            }
                        }
                    });
                },
                
                handleFileChange(e) {
                    const file = e.target.files[0];
                    if (file) {
                        this.formData.photo = file;
                        this.formData.photo_preview = URL.createObjectURL(file);
                    }
                },
                
                updateProfile() {
                    this.loading = true;
                    const token = localStorage.getItem('access_token');
                    const fd = new FormData();
                    fd.append('name', this.formData.name);
                    fd.append('email', this.formData.email);
                    if (this.formData.photo) fd.append('photo', this.formData.photo);
                    
                    fetch('/api/v1/profile', {
                        method: 'PUT',
                        headers: { 'Authorization': 'Bearer ' + token },
                        body: fd
                    })
                    .then(r => r.json())
                    .then(res => {
                        this.loading = false;
                        if (res.success) {
                            alert(res.message || 'Profil berhasil diperbarui');
                            this.fetchProfile();
                        } else {
                            alert('Gagal: ' + res.message);
                        }
                    })
                    .catch(() => { this.loading = false; alert('Terjadi kesalahan'); });
                },
                
                updatePassword() {
                    this.loadingPwd = true;
                    const token = localStorage.getItem('access_token');
                    fetch('/api/v1/profile/password', {
                        method: 'PUT',
                        headers: { 
                            'Authorization': 'Bearer ' + token,
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify(this.passwordData)
                    })
                    .then(r => r.json())
                    .then(res => {
                        this.loadingPwd = false;
                        if (res.success) {
                            alert('Password berhasil diubah');
                            this.passwordData = { old_password: '', new_password: '' };
                        } else {
                            alert('Gagal: ' + res.message);
                        }
                    })
                    .catch(() => { this.loadingPwd = false; alert('Terjadi kesalahan'); });
                }
            }));
        });
    </script>
"""
content = re.sub(r'<script>\s*document.addEventListener\(\'alpine:init\', \(\) => \{.*</script>', script_content, content, flags=re.DOTALL)
# also remove dashboardData reference
content = re.sub(r'x-data="dashboardData\(\)"', '', content)

with open('admin/profile.html', 'w') as f:
    f.write(content)
