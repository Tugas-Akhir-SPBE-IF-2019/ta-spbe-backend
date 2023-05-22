# this code is for TEXT FINDING - INDIKATOR 1
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
    init_keyword = list(keyword)

    # read line by line from txt
    for line in file:

        if (lv == 1):
            if re.search(keyword[0], line, re.IGNORECASE):
                result.append([idx, line])

        elif (lv == 2):
            # cek dgn keyword 'domain'. cont: domain arsitektur proses bisnis, domain data dan informasi, domain arsitektur keamanan, dll
            for key in init_keyword:
                if key in keyword:

                    # manfaatin regex untuk searching yg not only kata keyword
                    reg = f'(?:(domain)\s+(arsitektur)?\s?({key}))'
                    if re.search(reg, line, re.IGNORECASE):
                        result.append([idx, line])
                        keyword.remove(key)

        elif (lv == 3):
            # masih pengecekan keyword yg sama dgn level 2, tapi tanpa 'domain'
            for key in init_keyword:
                if key in keyword:

                    # gausah regex krn searching nya langsung kata keyword tanpa ada embel embel 'domain'
                    if re.search(key, line, re.IGNORECASE):
                        result.append([idx, line])

                        # hapus key yg udh ditemuin supaya gausah di search lagi
                        keyword.remove(key)

        else:  # lv == 4
            # cek dulu keyword utama level 4 which is arsitektur spbe nasional
            reg = f'(?:(arsitektur)\s+(spbe)\s+(nasional))'

            # kalau ketemu, lanjut cek keyword [integrasi, reviu, dll]
            if (re.search(reg, line, re.IGNORECASE)):
                # check if line contains any keyword from list
                res = [ele for ele in keyword if (ele in line)]
                if (res):
                    result.append([idx, line])

        idx += 1

    file.close()
    return (result)


def ceklvl(filename):
    list_final = []
    text_final = ''

    lvl1 = convert_keywords(["Arsitektur SPBE"])
    res1 = txtreader(filename, 1, lvl1)

    # cek if keyword lvl1 is not found, then return as empty string
    if (not res1):
        return text_final

    lvl2 = convert_keywords(["Proses Bisnis", "Data dan Informasi", "Infrastruktur SPBE",
            "Aplikasi SPBE", "Keamanan SPBE", "Layanan SPBE"])
    res2 = txtreader(filename, 2, lvl2)

    if (not res2):
        res2 = txtreader(filename, 3, lvl2)

    for el in res2:
        if (el[1] not in list_final):
            list_final.append(el[1])

    if (not res2):
        list_final.append(res1[0][1])
        return clean_text(list_final)

    # if lvl2 tdk terpenuhi (not all keyword found), end pengecekan lvl
    if (lvl2):
        return clean_text(list_final)

    lvl4 = convert_keywords(["integrasi", "reviu", "diselaraskan", "berpedoman", "perubahan"])
    res4 = txtreader(filename, 4, lvl4)

    if (not res4):
        return clean_text(list_final)

    for el in res4:
        if (el[1] not in list_final):
            list_final.append(el[1])

    return clean_text(list_final)
