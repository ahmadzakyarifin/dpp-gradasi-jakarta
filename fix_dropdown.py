import os
import re

html_files = [f for f in os.listdir('admin') if f.endswith('.html')]

for file in html_files:
    filepath = os.path.join('admin', file)
    with open(filepath, 'r') as f:
        content = f.read()
    
    # Remove `hidden` class and `:class="{'hidden': !userMenuOpen}"`
    # Also we can add `style="display: none;"`
    
    content = content.replace(
        'class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 py-1 z-50 hidden"\n                         :class="{\'hidden\': !userMenuOpen}"',
        'class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 py-1 z-50"\n                         style="display: none;"'
    )
    
    with open(filepath, 'w') as f:
        f.write(content)
