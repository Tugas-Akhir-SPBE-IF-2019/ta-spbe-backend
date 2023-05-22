# this code is for TEXT FINDING - INDIKATOR 7
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
from .utility import *


def txtreader(filename, lv, keyword):
    # func to search keyword in txt file

    # open txt file
    # additional hacks
    filename = filename.split("/")[-1]
    file = open(f'cleaned_{filename}.txt', 'r')

    idx = 0
    result = []

    # read line by line from txt
    for line in file:

        if (lv == 2):
            # check using 
            for key in keyword:
                reg = f'{key}'
                if re.search(reg, line, re.IGNORECASE):
                    result.append([idx, line])
        
        # lv 3 and 4 have the same logic to check for found text
        elif (lv == 3 or lv == 4):
            # check using main keyword first
            reg = f'(?:(sistem)\s*(penghubung)\s*(layanan))'
            if (re.search(reg, line, re.IGNORECASE)):
                for key in keyword:
                    if re.search(key, line, re.IGNORECASE):
                        result.append([idx, line])

        idx += 1

    file.close()
    return (result)


# if __name__ == '__main__':
#     pdfparser(sys.argv[1])


# for indicator 7, ceklvl will only check against lvl2 keywords only
# because the rest are related to lvl2 and will be checked in the next step (text similarity)
def ceklvl(filename):
    list_final = []

    lvl2 = convert_keywords(['sistem penghubung layanan'])
    res2 = txtreader(filename, 2, lvl2)

    # check if keyword lvl2 is not found, then return as empty string
    if (not res2):
        return ''

    for el in res2:
        list_final.append(el[1])
    
    lvl3 = convert_keywords([
        'seluruh opd', 
        'setiap opd'
        'seluruh unit kerja', 
        'setiap unit kerja'
        'seluruh pemerintah daerah',
        'setiap pemerintah daerah'])
    res3 = txtreader(filename, 3, lvl3)

    #immediately return if no result is found for lvl3
    if (not res3):
        return clean_text(list_final)

    for el in res3:
        if (el[1] not in list_final):
            list_final.append(el[1])

    lvl4 = convert_keywords(['keterhubungan', 
            'hubung', 
            'integrasi', 
            'berpedoman', 
            'reviu', 
            'diselaraskan', 
            'perubahan'])
    res4 = txtreader(filename, 4, lvl4)

    for el in res4:
        if (el[1] not in list_final):
            list_final.append(el[1])

    return clean_text(list_final)