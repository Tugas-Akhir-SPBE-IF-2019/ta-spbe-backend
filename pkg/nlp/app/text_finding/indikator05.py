# this code is for TEXT FINDING - INDIKATOR 5
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
from .utility import *

def txtreader(filename, lv, keyword):

    # additional hacks
    filename = filename.split("/")[-1]
    file = open(f'cleaned_{filename}.txt', 'r')

    idx = 0
    result = []
    
    reg1 = f'(?:(layanan)?\s*?((pusat)\s*(data)))'
    reg2 = f'(?:(pusat)?\s*?((data)\s*(center)))'

    for line in file:
        if (lv == 1):
            if (re.search(reg1, line.lower(), re.IGNORECASE) or re.search(reg2, line.lower(), re.IGNORECASE)):
                result.append([idx, line])

        else:
            if (re.search(reg1, line.lower(), re.IGNORECASE) or re.search(reg2, line.lower(), re.IGNORECASE)):
                for key in keyword:
                    if re.search(key, line.lower(), re.IGNORECASE):
                        result.append([idx, line])

        idx += 1

    file.close()
    return (result)

def ceklvl(filename):
    list_final = []

    lvl1 = []
    res1 = txtreader(filename, 1, lvl1)

    if (not res1):
        return ''

    lvl2 = convert_keywords([
        'organisasi perangkat daerah', 
        'opd', 
        'unit kerja', 
        'pemerintah daerah', 
        'perangkat daerah'])
    res2 = txtreader(filename, 2, lvl2)

    if (not res2):
        list_final.append(res1[0][1])
        return clean_text(list_final)

    for el in res2:
        if (el[1] not in list_final):
            list_final.append(el[1])
    
    lvl3 = convert_keywords(['setiap opd', 
           'seluruh opd', 
           'setiap unit kerja', 
           'seluruh unit kerja', 
           'setiap pemerintah daerah',
           'seluruh pemerintah daerah',
           'seluruh Perangkat Daerah'])
    res3 = txtreader(filename, 3, lvl3)

    if (not res3):
        return clean_text(list_final)

    for el in res3:
        if (el[1] not in list_final):
            list_final.append(el[1])

    lvl4 = convert_keywords(['hubung', 
            'sambung',
            'integrasi', 
            'berpedoman', 
            'reviu', 
            'diselaraskan', 
            'perubahan', 
            'interkoneksi', 
            'periodik'])
    res4 = txtreader(filename, 4, lvl4)

    for el in res4:
        if (el[1] not in list_final):
            list_final.append(el[1])

    return clean_text(list_final)