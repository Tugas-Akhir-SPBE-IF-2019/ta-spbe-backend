# Import libraries
import json
import pandas as pd
import numpy as np

import matplotlib.pyplot as plt
import seaborn as sns

from sklearn import metrics, manifold
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report, confusion_matrix,accuracy_score
from sklearn.feature_extraction.text import TfidfVectorizer

from sklearn.feature_extraction.text import TfidfTransformer
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.metrics.pairwise import cosine_similarity

import re

import transformers

import nltk
from nltk.tokenize import word_tokenize
from nltk.corpus import stopwords
from nltk.stem import WordNetLemmatizer
# Stemmer nltk
from nltk.stem import PorterStemmer

# Stemmer sastrawi
from Sastrawi.Stemmer.StemmerFactory import StemmerFactory

import ast
import math

import pickle
from sklearn import pipeline
from sklearn import linear_model
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.tree import DecisionTreeClassifier
from sklearn.metrics import classification_report
from sklearn.tree import export_graphviz
import graphviz

from Sastrawi.Stemmer.StemmerFactory import StemmerFactory

nltk.download('stopwords')
nltk.download('punkt')
nltk.download('averaged_perceptron_tagger')
nltk.download('wordnet')
nltk.download('omw-1.4')

# 'Pengelolaan data dilakukan melalui proses berikut: data master, data referensi, basis data, kualitas data, interoperabilitas data, dan arsitektur data.'
# kelola data proses data master data referensi basis data kualitas data interoperabilitas data arsitektur data
# data master data referensi basis data kualitas data interoperabilitas data arsitektur data
def feature_extraction_level_1(data):
  level_1 = []
  for i in range(len(data)):
      if ('draft' in data['Judul'][i] or 'draft' in data['teks'][i]):
          level_1.append(1)
      else:
          level_1.append(0)
  return level_1

def feature_extraction_level_2(data, level_1):
  level_2 = []
  for i in range(len(data)):
      if (level_1[i] == 0 and ('manajemen data' in data['teks'][i]  or 'kelola data' in data['teks'][i])):
          level_2.append(1)
      else:
          level_2.append(0)
  return level_2


def feature_extraction_level_3(data, level_1, level_2):
  level_3 = []
  for i in range(len(data)):
      if (level_1[i] == 0 and level_2[i] == 1 and
          (('arsitektur data' in data['teks'][i]) or
          ('data induk' in data['teks'][i] or 'data master' in data['teks'][i]) or 
          ('data referensi' in data['teks'][i]) or
          ('basis data' in data['teks'][i] ) or 
          ('kualitas data' in data['teks'][i]) or 
          ('interoperabilitas' in data['teks'][i] or 'interoperabilitas data' in data['teks'][i])
        )):
          level_3.append(1)
      else:
          level_3.append(0)
  return level_3

def feature_extraction_level_4(data, level_1, level_2, level_3):
  level_4 = []
  for i in range(len(data)):
      if (level_1[i] == 0 and level_2[i] == 1 and level_3[i] == 1 and
          (('arsitektur data' in data['teks'][i]) and
          ('data induk' in data['teks'][i] or 'data master' in data['teks'][i]) and
          ('data referensi' in data['teks'][i]) and
          ('basis data' in data['teks'][i] ) and 
          ('kualitas data' in data['teks'][i]) and 
          ('interoperabilitas' in data['teks'][i] or 'interoperabilitas data' in data['teks'][i])
        )):
          level_4.append(1)
      else:
          level_4.append(0)
  return level_4

# integrasi|reviu|diselaraskan|berpedoman|perubahan
# CEK

# arsitektur spbe|review|regulasi|ubah
# notes: level 5 masih salah ekstraksi fiturnya
def feature_extraction_level_5(data, level_1, level_2, level_3, level_4):
  level_5 = []
  for i in range(len(data)):
      level_5_per_item = [] 
      if level_1[i] == 0 and level_2[i] == 1 and level_3[i] == 1 and level_4[i] == 1:
        for j in range(len(data['JudulNotulensiUndangan'][i])):
          if ('sempurna' in data['JudulNotulensiUndangan'][i][j] or (('manajemen data' in data['IsiNotulensiUndangan'][i] or 'kelola data' in data['IsiNotulensiUndangan'][i]) or 'review' in data['IsiNotulensiUndangan'][i][j] or 'regulasi' in data['IsiNotulensiUndangan'][i][j] or 'ubah' in data['IsiNotulensiUndangan'][i][j])):
            level_5_per_item.append(1)
          else:
            level_5_per_item.append(0)

        if level_5_per_item == []:
          level_5.append(0)
        else:
          level_5.append(1)
      else:
          level_5.append(0)

  return level_5

def feature_extraction(df):
  level_1 = feature_extraction_level_1(df)
  level_2 = feature_extraction_level_2(df, level_1)
  level_3 = feature_extraction_level_3(df, level_1, level_2)
  level_4 = feature_extraction_level_4(df, level_1, level_2, level_3)
  level_5 = feature_extraction_level_5(df, level_1, level_2, level_3, level_4)

  df['level1_keyword'] = level_1
  df['level2_keyword'] = level_2
  df['level3_keyword'] = level_3
  df['level4_keyword'] = level_4
  df['level5_keyword'] = level_5
  
  return df

def convert_list(val):
    try:
        val = ast.literal_eval(val)
        if (val == ""):
            return None
        else:
            return val
    except (SyntaxError, ValueError):
        return val

def lowercase_sentence(text):
  text = text.lower()
  
  return text

def lowercase_sentence_in_list(text):
  for i in range(len(text)):
    text[i] = text[i].lower()
  
  return text

def repair_type(df):
  df['instansi'] = df['instansi'].astype(str)
  df['indikator'] = df['indikator'].astype(int)
  df['level'] = df['level'].astype(int)
  df['Judul'] = df['Judul'].astype(str)
  df['teks'] = df['teks'].astype(str)
  df['JudulDokumenLama'] = df['JudulDokumenLama'].apply(convert_list)
  df['TeksDokumenLama'] = df['TeksDokumenLama'].apply(convert_list)
  df['JudulNotulensiUndangan'] = df['JudulNotulensiUndangan'].apply(convert_list)
  df['IsiNotulensiUndangan'] = df['IsiNotulensiUndangan'].apply(convert_list)

def remove_between_square_brackets(text):
  return re.sub('\[[^]]*\]', '', text)

def remove_special_characters(text, remove_digits=True):
  pattern=r'[^a-zA-z0-9\s]'
  text=re.sub(pattern,'',text)
  return text

def remove_special_characters_in_list(text, remove_digits=True):
  for i in range(len(text)):
    pattern=r'[^a-zA-z0-9\s]'
    text[i]=re.sub(pattern,'',text[i])
  
  return text

def remove_stopwords(text, is_lower_case=False):
  tokens = word_tokenize(text)
  
  if is_lower_case:
    filtered_tokens = [token for token in tokens if token not in stopword_list]
  else:
    filtered_tokens = [token for token in tokens if token.lower() not in stopword_list]
  
  filtered_text = ' '.join(filtered_tokens)
  
  return filtered_text

stopword_list = stopwords.words('indonesian')
factory = StemmerFactory()
stemmer = factory.create_stemmer()

def remove_stopwords_in_list(text, is_lower_case=False):
  hasil = []
  for i in range(len(text)):
    tokens = word_tokenize(text[i])
  
    if is_lower_case:
      filtered_tokens = [token for token in tokens if token not in stopword_list]
    else:
      filtered_tokens = [token for token in tokens if token.lower() not in stopword_list]
    
    filtered_text = ' '.join(filtered_tokens)

    hasil.append(filtered_text)
  
  
  return hasil

def apply_preprocess(df):
  for col_name in df.columns[3:5]:
    df[col_name]=df[col_name].apply(lowercase_sentence)
    df[col_name]=df[col_name].apply(remove_between_square_brackets)
    df[col_name]=df[col_name].apply(remove_special_characters)
    df[col_name]=df[col_name].apply(remove_stopwords)
    df[col_name]=df[col_name].apply(lemmatization_sentence)

def apply_preprocess_list(df):
  for col_name in df.columns[5:]:
    df[col_name]=df[col_name].apply(lowercase_sentence_in_list)
    # df[col_name]=df[col_name].apply(remove_between_square_brackets)
    df[col_name]=df[col_name].apply(remove_special_characters_in_list)
    df[col_name]=df[col_name].apply(remove_stopwords_in_list)
    df[col_name]=df[col_name].apply(lemmatization_sentence_in_list)

def lemmatization_sentence(text):
  tokens = text.split()
  
  lemmatized_tokens = [stemmer.stem(word) for word in tokens]
  lemmatized_text = ' '.join(lemmatized_tokens)

  return lemmatized_text

def lemmatization_sentence_in_list(text):
  hasil = []
  for i in range(len(text)):
    tokens = text[i].split()
    
    lemmatized_tokens = [stemmer.stem(word) for word in tokens]
    lemmatized_text = ' '.join(lemmatized_tokens)

    hasil.append(lemmatized_text)

  return hasil

def preprocess_data(df):
  repair_type(df)
  apply_preprocess(df)
  apply_preprocess_list(df)
  feature_extraction(df)
  return df

def train_model(X_train, y_train):
  clf = DecisionTreeClassifier()
  clf.fit(X_train, y_train)
  return clf

def test_model(clf, X_test):
  y_pred = clf.predict(X_test)
  return y_pred

def repair_type(df):
  df['instansi'] = df['instansi'].astype(str)
  df['indikator'] = df['indikator'].astype(int)
  df['level'] = df['level'].astype(int)
  df['Judul'] = df['Judul'].astype(str)
  df['teks'] = df['teks'].astype(str)
  df['JudulDokumenLama'] = df['JudulDokumenLama'].apply(convert_list)
  df['TeksDokumenLama'] = df['TeksDokumenLama'].apply(convert_list)
  df['JudulNotulensiUndangan'] = df['JudulNotulensiUndangan'].apply(convert_list)
  df['IsiNotulensiUndangan'] = df['IsiNotulensiUndangan'].apply(convert_list)

def repair_type_for_model(df):
  df['Judul'] = df['Judul'].astype(str)
  df['teks'] = df['teks'].astype(str)
  df['JudulDokumenLama'] = df['JudulDokumenLama'].apply(convert_list)
  df['TeksDokumenLama'] = df['TeksDokumenLama'].apply(convert_list)
  df['JudulNotulensiUndangan'] = df['JudulNotulensiUndangan'].apply(convert_list)
  df['IsiNotulensiUndangan'] = df['IsiNotulensiUndangan'].apply(convert_list)

def apply_preprocess(df):
  for col_name in df.columns[3:5]:
    df[col_name]=df[col_name].apply(lowercase_sentence)
    df[col_name]=df[col_name].apply(remove_between_square_brackets)
    df[col_name]=df[col_name].apply(remove_special_characters)
    df[col_name]=df[col_name].apply(remove_stopwords)
    df[col_name]=df[col_name].apply(lemmatization_sentence)

def apply_preprocess_list(df):
  for col_name in df.columns[5:]:
    df[col_name]=df[col_name].apply(lowercase_sentence_in_list)
    # df[col_name]=df[col_name].apply(remove_between_square_brackets)
    df[col_name]=df[col_name].apply(remove_special_characters_in_list)
    df[col_name]=df[col_name].apply(remove_stopwords_in_list)
    df[col_name]=df[col_name].apply(lemmatization_sentence_in_list)


# for model
def apply_preprocess_for_model(df):
  for col_name in df.columns[:2]:
    df[col_name]=df[col_name].apply(lowercase_sentence)
    df[col_name]=df[col_name].apply(remove_between_square_brackets)
    df[col_name]=df[col_name].apply(remove_special_characters)
    df[col_name]=df[col_name].apply(remove_stopwords)
    df[col_name]=df[col_name].apply(lemmatization_sentence)

def apply_preprocess_for_model_list(df):
  for col_name in df.columns[2:]:
    df[col_name]=df[col_name].apply(lowercase_sentence_in_list)
    # df[col_name]=df[col_name].apply(remove_between_square_brackets)
    df[col_name]=df[col_name].apply(remove_special_characters_in_list)
    df[col_name]=df[col_name].apply(remove_stopwords_in_list)
    df[col_name]=df[col_name].apply(lemmatization_sentence_in_list)

def predict_model(clf, df):
  repair_type_for_model(df)
  apply_preprocess_for_model(df)
  apply_preprocess_for_model_list(df)
  feature_extraction(df)
  print(df)
  df = df.drop(['Judul', 'JudulDokumenLama', 'TeksDokumenLama', 'teks', 'JudulNotulensiUndangan', 'IsiNotulensiUndangan'], axis=1)
  df = clf.predict(df)
  return df

# def print_tree(clf):
#   dot_data = export_graphviz(
#       clf, 
#       out_file=None, 
#       feature_names=X.columns,
#       class_names=["lvl 1", "lvl 2", "lvl 3", "lvl 4", "lvl 5"],  
#       filled=True, rounded=True,  
#       special_characters=True)
  
#   graph = graphviz.Source(dot_data)
#   graph.render("spbe")