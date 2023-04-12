# this code is for TEXT FINDING - INDIKATOR 4
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
import preprocess_dokbaru as dokbaru
import preprocess_doklama as doklama
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
    exclude_keywords = ['badan', 'anggaran', 'belanja', 'TIK', 'manajemen', 'fungsi',
               'infrastruktur', 'interoperabilitas', 'rencana induk', 'program', 'pelayanan']

    # read line by line from txt
    for line in file:

        if (lv == 1):
            reg = f'(?:(pembangunan)(.*?)(aplikasi)\s+(spbe)?)'
            if (re.search(reg, line, re.IGNORECASE)):
                result.append([idx, line])

        elif (lv == 2):
            for key in init_keyword:  # create a duplicate list
                if key in keyword:  # check if element in iteration exist in otiginal list, skip element if not exist
                    if re.search(key, line, re.IGNORECASE):

                        # to include ['aplikasi', 'aplikasi khusus', 'aplikasi umum', 'aplikasi spbe']
                        reg = f'(?:(aplikasi)\s+[(umum)|(khusus)|(spbe)]?)'
                        if (re.search(reg, line, re.IGNORECASE) and not exclude_words(exclude_keywords, line)):
                            result.append([idx, line])
                            # remove element in original list
                            keyword.remove(key)

        elif (lv == 3):
            for key in init_keyword:
                if key in keyword:
                    if re.search(key, line, re.IGNORECASE):
                        reg = f'(?:(aplikasi)\s+[(umum)|(khusus)|(spbe)]?)'
                        if re.search(reg, line, re.IGNORECASE):
                            result.append([idx, line])
                            keyword.remove(key)

        elif (lv == 4):
            res = [ele for ele in keyword if (ele in line)]
            if (len(res) == len(keyword)):  # if all keyword is found in the same sentence
                reg = f'(?:(aplikasi)\s+[(umum)|(khusus)|(spbe)]?)'
                if (re.search(reg, line, re.IGNORECASE) and not exclude_words(exclude_keywords, line)):
                    result.append([idx, line])

        else:
            reg = f'(?:(aplikasi)\s+(spbe))'
            if (re.search(reg, line.lower(), re.IGNORECASE)):
                if (re.search(keyword[0], line.lower(), re.IGNORECASE) and not exclude_words(exclude_keywords, line)):
                    result.append([idx, line])

        idx += 1

    file.close()
    return (result)


def ceklvl(filename):
    list_final = []
    text_final = ''

    res1 = txtreader(filename, 1, [])

    # cek if keyword lvl1 is not found, then return as empty string
    if (not res1):
        return text_final

    lvl2 = ["perencanaan", "analisis", "desain", "implementasi",
            "pemeliharaan"]
    res2 = txtreader(filename, 2, lvl2)

    if (not res2):
        list_final.append(res1[0][1])
        return clean_text(list_final)

    for el in res2:
        if (el[1] not in list_final):
            list_final.append(el[1])

    # if lvl2 tdk terpenuhi (not all keyword found), end pengecekan lvl
    if (lvl2):
        return clean_text(list_final)

    lvl3 = ["konsultasi", "koordinasi"]
    res3 = txtreader(filename, 3, lvl3)

    if (not res3):
        return clean_text(list_final)

    for el in res3:
        if (el[1] not in list_final):
            list_final.append(el[1])

    lvl4 = ["terpadu", "endali"]
    res4 = txtreader(filename, 4, lvl4)

    if (not res4):
        return clean_text(list_final)

    for el in res4:
        if (el[1] not in list_final):
            list_final.append(el[1])

    lvl5 = ["reviu"]
    res5 = txtreader(filename, 5, lvl4)

    for el in res5:
        if (el[1] not in list_final):
            list_final.append(el[1])

    return clean_text(list_final)