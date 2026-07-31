import re

with open('admin/users.html', 'r') as f:
    content = f.read()

# Remove column header
content = re.sub(r'<th class="p-4">Kontak \(No HP\)</th>\s*', '', content)
# Remove column data
content = re.sub(r'<td class="p-4">\s*<div class="flex items-center gap-2">\s*<i class="ph ph-phone text-gray-400"></i>\s*<span class="text-gray-600" x-text="item.phone \|\| \'-\'"></span>\s*</div>\s*</td>', '', content)
# It appears twice (for different status branches perhaps)
content = re.sub(r'<td class="p-4">\s*<div class="flex items-center gap-2">\s*<i class="ph ph-phone text-gray-400"></i>\s*<span class="text-gray-600" x-text="item\.phone \|\| \'-\'"></span>\s*</div>\s*</td>', '', content)

# Remove input field
input_pattern = r'<div>\s*<label class="block text-sm font-medium text-gray-700 mb-1">No Handphone</label>\s*<input type="tel" x-model="formData\.phone" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm">\s*</div>\s*'
content = re.sub(input_pattern, '', content)

with open('admin/users.html', 'w') as f:
    f.write(content)

