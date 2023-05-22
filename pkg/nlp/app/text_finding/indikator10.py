# this code is for TEXT FINDING - INDIKATOR 10
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
    exclude_keywords = ['badan', 'anggaran', 'belanja', 'TIK', 'manajemen', 'fungsi',
               'infrastruktur', 'interoperabilitas', 'rencana induk', 'program', 'pelayanan']
    init_keywords = list(keyword)

    # read line by line from txt
    for line in file:

        if (lv == 1):
            reg = f'(?:(tim)\s+(koordinasi)\s+(spbe)?)'
            if (re.search(reg, line, re.IGNORECASE)):
                result.append([idx, line])

        elif (lv == 2):
            for key in init_keywords:
                if key in keyword:
                    if re.search(key, line, re.IGNORECASE):

                        # include words
                        reg = f'(?:(tim)\s+(koordinasi)\s+(spbe)?)'
                        if (re.search(reg, line, re.IGNORECASE) and not exclude_words(exclude_keywords, line)):
                            result.append([idx, line])
                            keyword.remove(key)

        elif (lv == 4):
            reg = f'(?:(tim)\s+(koordinasi))'
            if re.search(reg, line, re.IGNORECASE):
                res = [ele for ele in keyword if (ele in line)]
                if (res):
                    if (not exclude_words(exclude_keywords, line)):
                        result.append([idx, line])

        else:  # lvl 5
            res = [ele for ele in keyword if (ele in line)]
            if (res):
                reg = f'(?:(tim)\s+(koordinasi)\s+(spbe)?)'
                if re.search(reg, line, re.IGNORECASE):
                    if (not exclude_words(exclude_keywords, line)):
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
        return ''

    # TODO
    # add setiap opd, etc

    lvl2 = ['tugas']
    res2 = txtreader(filename, 2, lvl2)

    if (not res2):
        list_final.append(res1[0][1])
        return clean_text(list_final)

    for el in res2:
        if (el[1] not in list_final):
            list_final.append(el[1])

    lvl4 = convert_keywords(["integrasi", "diselaraskan", "berpedoman",
            "perubahan", "koordinasi", "kerja sama"])
    res4 = txtreader(filename, 4, lvl4)

    if (not res4):
        return clean_text(list_final)

    for el in res4:
        if (el[1] not in list_final):
            list_final.append(el[1])

    lvl5 = ["reviu"]
    res5 = txtreader(filename, 5, [])

    for el in res5:
        if (el[1] not in list_final):
            list_final.append(el[1])

    return clean_text(list_final)