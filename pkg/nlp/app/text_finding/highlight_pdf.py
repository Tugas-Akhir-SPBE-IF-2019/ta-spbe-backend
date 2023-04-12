# Import Libraries
from typing import Tuple, List
from io import BytesIO
import os
import argparse
import re
import fitz
from .utility import *


def extract_info(input_file: str):
    """
    Extracts file info
    """
    # Open the PDF
    pdfDoc = fitz.open(input_file)
    output = {
        "File": input_file, "Encrypted": ("True" if pdfDoc.isEncrypted else "False")
    }
    # If PDF is encrypted the file metadata cannot be extracted
    if not pdfDoc.isEncrypted:
        for key, value in pdfDoc.metadata.items():
            output[key] = value

    # To Display File Info
    print("## File Information ##################################################")
    print("\n".join("{}:{}".format(i, j) for i, j in output.items()))
    print("######################################################################")

    return True, output


def search_for_text(lines, search_str):
    """
    Search for the search string within the document lines
    """
    for line in lines:
        # Find all matches within one line
        results = re.findall(search_str, line, re.IGNORECASE)
        # In case multiple matches within one line
        for result in results:
            yield result


def redact_matching_data(page, matched_values):
    """
    Redacts matching values
    """
    matches_found = 0
    # Loop throughout matching values
    for val in matched_values:
        matches_found += 1
        matching_val_area = page.search_for(val)
        # Redact matching values
        [page.addRedactAnnot(area, text=" ", fill=(0, 0, 0))
         for area in matching_val_area]
    # Apply the redaction
    page.apply_redactions()
    return matches_found


def frame_matching_data(page, matched_values):
    """
    frames matching values
    """
    matches_found = 0
    # Loop throughout matching values
    for val in matched_values:
        matches_found += 1
        matching_val_area = page.search_for(val)
        for area in matching_val_area:
            if isinstance(area, fitz.fitz.Rect):
                # Draw a rectangle around matched values
                annot = page.addRectAnnot(area)
                # , fill = fitz.utils.getColor('black')
                annot.setColors(stroke=fitz.utils.getColor('red'))
                # If you want to remove matched data
                #page.addFreetextAnnot(area, ' ')
                annot.update()
    return matches_found


def highlight_matching_data(page, matched_values, type, matched_ends=None):
    """
    Highlight matching values
    """
    matches_found = 0
    # Loop throughout matching values
    for val in matched_values:
        matches_found += 1
        matching_val_area = page.search_for(val)
        highlight = None
        if type == 'Squiggly':
            highlight = page.addSquigglyAnnot(matching_val_area)
        elif type == 'Underline':
            highlight = page.addUnderlineAnnot(matching_val_area)
        elif type == 'Strikeout':
            highlight = page.addStrikeoutAnnot(matching_val_area)
        else:
            if (matched_ends == None):
                highlight = page.add_highlight_annot(matching_val_area)
            else:
                for end in matched_ends:
                    matching_end_area = page.search_for(end)
                    highlight = page.add_highlight_annot(
                        start=matching_val_area[0].tl, stop=matching_end_area[0].br)
        # To change the highlight colar
        # highlight.setColors({"stroke":(0,0,1),"fill":(0.75,0.8,0.95) })
        # highlight.setColors(stroke = fitz.utils.getColor('white'), fill = fitz.utils.getColor('red'))
        # highlight.setColors(colors= fitz.utils.getColor('red'))
        if (highlight is not None):
            highlight.update()
    return matches_found


def process_data(input_file: str, output_file: str, search_str, pages: Tuple = None, action: str = 'Highlight'):
    """
    Process the pages of the PDF File
    """
    # Open the PDF
    pdfDoc = fitz.open(input_file)
    # Save the generated PDF to memory buffer
    output_buffer = BytesIO()
    page_with_matches = []
    total_matches = 0
    # Iterate through pages
    for pg in range(pdfDoc.page_count):
        # If required for specific pages
        if pages:
            if str(pg) not in pages:
                continue
        # Select the page
        page = pdfDoc[pg]
        # Get Matching Data
        # Split page by lines
        page_lines = page.get_text("text").split('\n')
        if type(search_str) == str:
            matched_values = search_for_text(page_lines, search_str)
        else:
            matched_values = search_for_text(page_lines, search_str[0])
            matched_ends = search_for_text(page_lines, search_str[1])

            if not matched_ends:
                temp_page = pdfDoc[pg + 1]
                page_lines = page.get_text("text").split('\n')
                matched_ends = search_for_text(page_lines, search_str[1])

        if matched_values:
            if action == 'Redact':
                matches_found = redact_matching_data(page, matched_values)
            elif action == 'Frame':
                matches_found = frame_matching_data(page, matched_values)
            elif action in ('Highlight', 'Squiggly', 'Underline', 'Strikeout'):
                if type(search_str) == str:
                    matches_found = highlight_matching_data(
                        page, matched_values, action)
                else:
                    matches_found = highlight_matching_data(
                        page, matched_values, action, matched_ends)
            else:
                if type(search_str) == str:
                    matches_found = highlight_matching_data(
                        page, matched_values, 'Highlight')
                else:
                    matches_found = highlight_matching_data(
                        page, matched_values, 'Highlight', matched_ends)
            if matches_found > 0:
                page_with_matches.append(pg + 1)
            total_matches += matches_found
    print(f"{total_matches} Match(es) Found of Search String {search_str} In Input File: {input_file}")
    # Save to output
    pdfDoc.save(output_buffer)
    pdfDoc.close()
    # Save the output buffer to the output file
    with open(output_file, mode='wb') as f:
        f.write(output_buffer.getbuffer())

    return total_matches, page_with_matches


def remove_highlght(input_file: str, output_file: str, pages: Tuple = None):
    # Open the PDF
    pdfDoc = fitz.open(input_file)
    # Save the generated PDF to memory buffer
    output_buffer = BytesIO()
    # Initialize a counter for annotations
    annot_found = 0
    # Iterate through pages
    for pg in range(pdfDoc.page_count):
        # If required for specific pages
        if pages:
            if str(pg) not in pages:
                continue
        # Select the page
        page = pdfDoc[pg]
        annot = page.first_annot
        while annot:
            annot_found += 1
            page.delete_annot(annot)
            annot = annot.next
    if annot_found >= 0:
        print(f"Annotation(s) Found In The Input File: {input_file}")
    # Save to output
    pdfDoc.save(output_buffer)
    pdfDoc.close()
    # Save the output buffer to the output file
    with open(output_file, mode='wb') as f:
        f.write(output_buffer.getbuffer())


def process_file(**kwargs):
    """
    To process one single file
    Redact, Frame, Highlight... one PDF File
    Remove Highlights from a single PDF File
    """
    input_file = kwargs.get('input_file')
    output_file = kwargs.get('output_file')
    total_matches = 0
    page_with_matches = []
    if output_file is None:
        output_file = input_file
    search_str = kwargs.get('search_str')
    pages = kwargs.get('pages')
    # Redact, Frame, Highlight, Squiggly, Underline, Strikeout, Remove
    action = kwargs.get('action')
    if action == "Remove":
        # Remove the Highlights except Redactions
        remove_highlght(input_file=input_file,
                        output_file=output_file, pages=pages)
    else:
        total_matches, page_match = process_data(input_file=input_file, output_file=output_file,
                                     search_str=search_str, pages=pages, action=action)
        page_with_matches.extend(page_match)
    return total_matches, page_with_matches


def process_folder(**kwargs):
    """
    Redact, Frame, Highlight... all PDF Files within a specified path
    Remove Highlights from all PDF Files within a specified path
    """
    input_folder = kwargs.get('input_folder')
    search_str = kwargs.get('search_str')
    # Run in recursive mode
    recursive = kwargs.get('recursive')
    #Redact, Frame, Highlight, Squiggly, Underline, Strikeout, Remove
    action = kwargs.get('action')
    pages = kwargs.get('pages')
    # Loop though the files within the input folder.
    for foldername, dirs, filenames in os.walk(input_folder):
        for filename in filenames:
            # Check if pdf file
            if not filename.endswith('.pdf'):
                continue
             # PDF File found
            inp_pdf_file = os.path.join(foldername, filename)
            print("Processing file =", inp_pdf_file)
            process_file(input_file=inp_pdf_file, output_file=None,
                         search_str=search_str, action=action, pages=pages)
        if not recursive:
            break


def is_valid_path(path):
    """
    Validates the path inputted and checks whether it is a file path or a folder path
    """
    if not path:
        raise ValueError(f"Invalid Path")
    if os.path.isfile(path):
        return path
    elif os.path.isdir(path):
        return path
    else:
        raise ValueError(f"Invalid Path {path}")

# Main code


def highlight(filename, searchstr, action):
    # Extracting File Info
    extract_info(input_file=filename)

    # Split search string by sentences
    reslist = re.split('[!?.,:;]', searchstr)
    page_with_matches = []

    # Clean res list
    for item in reslist:
        if len(item.strip()) <= 1 or item.isdigit():
            reslist.remove(item)

    # reglist = convert_search(reslist)
    # print(reglist)

    # Process a file
    for item in list(reslist):
        if (item in reslist):
            total_matches, page_matchs = process_file(
                input_file=filename,
                search_str=convert_search(item),
                action=action)
            page_with_matches.extend(page_matchs)

            if total_matches != 0:
                reslist.remove(item)

    for item in list(reslist):
        if (item in reslist):
            items = item.split()

            # Get initial word
            init_word = items[0] + '\s+' + items[1]
            for i in range(len(items)):
                if len(items[i].strip()) > 1 or not items[i].isdigit():
                    init_word = items[i] + '\s+' + items[i+1]
                    break
            # Get last word
            last_word = items[-1]

            total_matches, page_matchs = process_file(
                input_file=filename,
                search_str=[init_word, last_word],
                action=action)
            page_with_matches.extend(page_matchs)
    
    page_with_matches = list(dict.fromkeys(page_with_matches))
    pdfDoc = fitz.open(filename)
    proof_pic_files = []
    clean_filenames = []
    for page in (page_with_matches):
        pix = pdfDoc[page - 1].get_pixmap()
        pic_filename = f"{filename}-page-{page}.png"
        pix.save(pic_filename)

        clean_filename = filename.split("/")[-1]
        clean_pic_filename = f"{clean_filename}-page-{page}.png"
        proof_pic_files.append(clean_pic_filename)
        clean_filenames.append(clean_filename+f"#page={page}")
    return proof_pic_files, clean_filenames, page_with_matches

# filename = 'F2201-287-Indikator_01_+_Indikator1_Perbup_81_tahun_2021.pdf'
# test = "domain Arsitektur Proses Bisnis, b. domain Arsitektur data dan informasi, c. domain Arsitektur Infrastruktur SPBE, d. domain Arsitektur Aplikasi SPBE, e. domain Arsitektur Keamanan SPBE, f. domain Arsitektur Layanan SPBE, (7)  Untuk  menyelaraskan  Arsitektur  SPBE  Pemerintah Daerah  dengan  Arsitektur  SPBE  Nasional,  Bupati berkoordinasi  dan  dapat  melakukan  konsultasi  dengan menteri yang menyelenggarakan urusan pemerintahan di bidang aparatur Negara"
# page_with_matches = highlight(filename, test, 'Highlight')
# pdfDoc = fitz.open(filename)
# for page in (page_with_matches):
#     pix = pdfDoc[page - 1].get_pixmap()
#     pix.save("page-%i.png" % page)
