import re
import sys

def update_file(filename):
    with open(filename, 'r') as f:
        content = f.read()

    # Pattern for kegiatan.html
    # <p class="text-brand-600 text-[10px] font-bold tracking-wider uppercase mb-2 flex items-center gap-1.5">
    #     <i class="ph-bold ph-calendar-blank text-sm"></i> 31 Desember 2025
    # </p>
    
    # We will replace it with:
    # <div class="flex items-center gap-3 mb-2 flex-wrap">
    #     <p class="text-brand-600 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5">
    #         <i class="ph-bold ph-calendar-blank text-sm"></i> \1
    #     </p>
    #     <p class="text-slate-500 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 border-l border-slate-200 pl-3">
    #         <i class="ph-bold ph-map-pin text-sm"></i> \2
    #     </p>
    # </div>
    
    # For now, let's just make sure we capture the date and inject a default location if not present.
    # Actually, a simpler replace:
    
    def repl_keg(m):
        date = m.group(1).strip()
        # Find some location based on the date or just default
        loc = "Jakarta"
        if "November" in date: loc = "Bandung"
        if "Oktober" in date: loc = "Yogyakarta"
        if "Agustus" in date: loc = "Surabaya"
        
        return f'''<div class="flex items-center gap-3 mb-2 flex-wrap">
                            <p class="text-brand-600 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5">
                                <i class="ph-bold ph-calendar-blank text-sm"></i> {date}
                            </p>
                            <p class="text-slate-500 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 border-l border-slate-200 pl-3">
                                <i class="ph-bold ph-map-pin text-sm"></i> {loc}
                            </p>
                        </div>'''
                        
    pattern_keg = r'<p class="text-brand-600 text-\[10px\] font-bold tracking-wider uppercase mb-2 flex items-center gap-1\.5">\s*<i class="ph-bold ph-calendar-blank text-sm"></i>([^<]+)</p>'
    content = re.sub(pattern_keg, repl_keg, content)

    # Pattern for index.html
    # <p class="text-brand-600 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5">
    #     <i class="ph-bold ph-calendar-blank text-brand-500"></i> 31 December 2025
    # </p>
    
    def repl_idx(m):
        date = m.group(1).strip()
        loc = "Jakarta"
        return f'''<div class="flex items-center gap-3 flex-wrap">
                                <p class="text-brand-600 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5">
                                    <i class="ph-bold ph-calendar-blank text-brand-500"></i> {date}
                                </p>
                                <p class="text-slate-500 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 border-l border-slate-200 pl-3">
                                    <i class="ph-bold ph-map-pin text-slate-400"></i> {loc}
                                </p>
                            </div>'''

    pattern_idx = r'<p class="text-brand-600 text-\[10px\] font-bold tracking-wider uppercase flex items-center gap-1\.5">\s*<i class="ph-bold ph-calendar-blank text-brand-500"></i>([^<]+)</p>'
    content = re.sub(pattern_idx, repl_idx, content)

    with open(filename, 'w') as f:
        f.write(content)

update_file('kegiatan.html')
update_file('index.html')
