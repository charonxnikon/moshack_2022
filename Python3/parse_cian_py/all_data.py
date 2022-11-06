! pip3 install -U selenium
! pip3 install webdriver-manager
! pip3 install pandas 
! pip3 install requests
! pip3 install numpy
! pip3 install bs4
! pip3 install html
! pip3 install tqdm
! pip3 install gspread
! pip3 install df2gspread
! pip3 install time
! pip3 install chromedriver_autoinstaller
! pip3 install fake_useragent
! pip3 install datetime
! pip3 install undetected_chromedriver
! pip3 install webdriver_manager
! pip3 install warnings
! pip3 install geopy

import pandas as pd
import requests
import numpy as np
from bs4 import BeautifulSoup
import html
from tqdm import tqdm
import os
import gspread
import df2gspread as d2g
import time
import warnings
warnings.filterwarnings("ignore")

import chromedriver_autoinstaller as chromedriver
chromedriver.install()

!pip install --upgrade --force-reinstall chromedriver-binary-auto
import chromedriver_binary

from fake_useragent import UserAgent
import datetime
from selenium import webdriver
from xvfbwrapper import Xvfb
from selenium.webdriver.chrome.options import Options

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
import undetected_chromedriver

options = webdriver.ChromeOptions()
options.add_argument("user-agent=Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0")
options.add_argument("--disable-blink-features=AutomationControlled")

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
driver.get("https://www.google.com")
driver.close()

def get_html(url):
    req = requests.get(url, headers={"User-Agent": UserAgent().chrome})
    return req.text.encode(req.encoding)

def get_title(soup):
    try:
        title = soup.find("h1").text.strip()
    except Exception as e:
        #print(str(e) + " title")
        title = "Не указано"
    return title


def get_address(soup):
    try:
        address = soup.find("address").text.strip()
        if "На карте" in address:
            address = address[:address.rfind("На карте")]
        # separating data from the address string
        district, street = "Не указано", "Не указано"
        city = address.split(",")[1].strip()
        block_number = address.split(",")[-1].strip()
        if "ул " in block_number.lower() or "ул." in block_number.lower() or "улица" in block_number.lower() \
                or " пер" in block_number.lower() or "проезд" in block_number.lower() or "проспект" in block_number.lower():
            street = block_number
            block_number = "Не указано"

        for param in address.split(",")[1:-1]:
            if "ул " in param.lower() or "ул." in param.lower() or "улица" in param.lower() or " пер" in param.lower() \
                    or "проезд" in param.lower() or "проспект" in param.lower():
                street = param.strip()
            elif "район" in param.lower() or "р-н" in param.lower():
                district = param.strip()

        if street.split()[-1].strip().isdigit():
            block_number = street.split()[-1].strip()
            street = " ".join(street.split()[:-1]).strip()

        return city, district, street, block_number
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_title\n")
    return ["Не указано"] * 4


def get_price(soup):
    try:
        price = soup.find("span", {"itemprop": "price"})
        if price is not None:
            price = price.text.strip()
        else:
            price = "от " + soup.find("span", {"itemprop": "lowPrice"}).text.strip() + \
                    " до " + soup.find("span", {"itemprop": "highPrice"}).text.strip() + "/мес."
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_price\n")
        price = "Не указано"
    return price


def get_selling_type(soup):
    try:
        paragraphs = [x for x in soup.find_all("p") if x.get("class") is not None
                      and len(x.get("class")) == 1 and "description--" in x.get("class")[0]]
        if paragraphs:
            selling_type = paragraphs[0].text.strip()
        else:
            selling_type = "Не указано"
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_selling_type\n")
        selling_type = "Не указано"
    return selling_type


def get_seller_type(soup):
    try:
        divs = [x for x in soup.find_all("div") if x.get("class") is not None
                and len(x.get("class")) == 1 and "honest-container" in x.get("class")[0]]
        if not divs:
            seller_type = "Не указано"
        else:
            seller_type = divs[0].text.strip()
            if seller_type is not None and seller_type.lower() == "собственник":
                seller_type = "Собственник"
            else:
                seller_type = "Посредник"
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_seller_type\n")
        seller_type = "Не указано"
    return seller_type


def get_seller_name(soup):
    try:
        name = [x for x in soup.find_all("h2") if x.get("class") is not None and len(x.get("class")) == 1
                and "title--" in x.get("class")[0]]
        if name:
            name = name[0].text.strip()
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_seller_name\n")
        name = "Не указано"
    return name


def get_photos(url):
    try:
        driver = webdriver.Chrome()
        driver.get(url)

        images = []
        images_list = driver.find_elements_by_class_name("fotorama__img")
        images_list = [x.get_attribute("src") for x in images_list if "-2." in x.get_attribute("src")]
        for image in images_list:
            link = image.replace("-2.", "-1.")
            images.append(link)
        images = "\n".join(images)
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_photos\n")
        images = "Не указано"
    return images


def get_description(soup):
    try:
        paragraphs = [x for x in soup.find_all("p") if x.get("class") is not None
                      and len(x.get("class")) == 1 and "description-text--" in x.get("class")[0]]
        description = paragraphs[0].text.strip()
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_description\n")
        description = "Не указано"
    return description


def get_date(soup):
    try:
        date = soup.find("div", id="frontend-offer-card").find("main").find_all("div")[4].text.strip()
        if "вчера" in date:
            date = str(datetime.datetime.today() - datetime.timedelta(days=1)).split()[0]
        elif "сегодня" in date:
            date = str(datetime.datetime.today()).split()[0]
        else:
            date = "too old"
    except Exception as e:
        with open("logs.txt", "a", encoding="utf8") as file:
            file.write(str(e) + " cian get_date\n")
        date = "Не указано"
    return date

def find_kitchen(soup):
    temp_1 = soup.find_all('div', 'a10a3f92e9--info-value--bm3DC')
    temp_2 = soup.find_all('div','a10a3f92e9--info-title--JWtIm')
    j = 0
    flag = 0
    for i in temp_2:
        if i.text == 'Кухня':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_1[ind].text)

def find_floor(soup):
    temp_1 = soup.find_all('div', 'a10a3f92e9--info-value--bm3DC')
    temp_2 = soup.find_all('div','a10a3f92e9--info-title--JWtIm')
    j = 0
    flag = 0
    for i in temp_2:
        if i.text == 'Этаж':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_1[ind].text)

def find_finishing(soup):
    temp_1 = soup.find_all('div', 'a10a3f92e9--info-value--bm3DC')
    temp_2 = soup.find_all('div','a10a3f92e9--info-title--JWtIm')
    j = 0
    flag = 0
    for i in temp_2:
        if i.text == 'Отделка':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_1[ind].text)

def find_living(soup):
    temp_1 = soup.find_all('div', 'a10a3f92e9--info-value--bm3DC')
    temp_2 = soup.find_all('div','a10a3f92e9--info-title--JWtIm')
    j = 0
    flag = 0
    for i in temp_2:
        if i.text == 'Жилая':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_1[ind].text)

def find_year(soup):
    temp_1 = soup.find_all('div', 'a10a3f92e9--info-value--bm3DC')
    temp_2 = soup.find_all('div','a10a3f92e9--info-title--JWtIm')
    j = 0
    flag = 0
    for i in temp_2:
        if i.text == 'Построен':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_1[ind].text)

def type_floar(soup):
    temp_3 = soup.find_all('span', 'a10a3f92e9--name--x7_lt')
    temp_4 = soup.find_all('span','a10a3f92e9--value--Y34zN')
    j = 0
    flag = 0
    for i in temp_3:
        if i.text == 'Тип жилья':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_4[ind].text)
def type_san(soup):
    temp_3 = soup.find_all('span', 'a10a3f92e9--name--x7_lt')
    temp_4 = soup.find_all('span','a10a3f92e9--value--Y34zN')
    j = 0
    flag = 0
    for i in temp_3:
        if i.text == 'Санузел':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_4[ind].text)

def type_bal(soup):
    temp_3 = soup.find_all('span', 'a10a3f92e9--name--x7_lt')
    temp_4 = soup.find_all('span','a10a3f92e9--value--Y34zN')
    j = 0
    flag = 0
    for i in temp_3:
        if i.text == 'Балкон/лоджия' or i.text == 'Балкон':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_4[ind].text)
def type_rem(soup):
    temp_3 = soup.find_all('span', 'a10a3f92e9--name--x7_lt')
    temp_4 = soup.find_all('span','a10a3f92e9--value--Y34zN')
    j = 0
    flag = 0
    for i in temp_3:
        if i.text == 'Ремонт':
            ind = j
            flag = 1
        j += 1
    if flag == 0:
        return('Не указано')
    else:
        return(temp_4[ind].text)

def find_distance(soup):
    temp_5 = soup.find_all('span', 'a10a3f92e9--highway_distance--wUNBn')
    if len(temp_5) == 0:
        temp_5 = soup.find_all('span', 'a10a3f92e9--underground_time--iOoHy')
        temp_5_2 = soup.find_all('a', 'a10a3f92e9--underground_link--Sxo7K')
        if len(temp_5_2) == 0 or len(temp_5) == 0:
            return np.nan
        else: 
            if (temp_5_2[0].text + temp_5[0].text).find('МКАД') != -1:
                return np.nan
            else:
                return float((temp_5_2[0].text + temp_5[0].text).split('⋅  ')[1].split('мин')[0])
    else:
        return np.nan

from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException
#C:\\WebDrivers\\chromedriver.exe' для винды
driver = undetected_chromedriver.Chrome()
driver.get('https://www.cian.ru/cat.php?deal_type=sale&engine_version=2&offer_type=flat&region=1&room1=1&room2=1&room3=1&room4=1&room5=1&room6=1&sort=creation_date_desc')
page_source = driver.page_source
soup = BeautifulSoup(page_source, 'lxml')
temp_7 = soup.find_all('a', '_93444fe79c--link--eoxce')
to_iter = int(temp_7[0]['href'].split('/')[-2])
#options = webdriver.ChromeOptions()
#options.add_argument("start-maximized")
#options.add_experimental_option("excludeSwitches", ["enable-automation"])
#options.add_experimental_option('useAutomationExtension', False)
#driver = webdriver.Chrome(chrome_options=options,executable_path='/Users/Nikon/.wdm/drivers/chromedriver/mac64/107.0.5304.62/chromedriver')
#driver.get("https://sslproxies.org/")
#driver = webdriver.Chrome()
#driver.execute_script("return arguments[0].scrollIntoView(true);", WebDriverWait(driver, 20).until(EC.visibility_of_element_located((By.XPATH, "//table[@class='table table-striped table-bordered dataTable']//th[contains(., 'IP Address')]"))))
#ips = [my_elem.get_attribute("innerHTML") for my_elem in WebDriverWait(driver, 5).until(EC.visibility_of_all_elements_located((By.XPATH, "//table[@class='table table-striped table-bordered dataTable']//tbody//tr[@role='row']/td[position() = 1]")))]
#ports = [my_elem.get_attribute("innerHTML") for my_elem in WebDriverWait(driver, 5).until(EC.visibility_of_all_elements_located((By.XPATH, "//table[@class='table table-striped table-bordered dataTable']//tbody//tr[@role='row']/td[position() = 2]")))]
#driver.quit()
driver.close()
driver = undetected_chromedriver.Chrome()
resp = requests.get('https://free-proxy-list.net/') 
df = pd.read_html(resp.text)[0]
proxies = []
for i in range(0,len(df)):
    proxies.append(str(df['IP Address'][i])+':'+str(df['Port'][i]))
j = -1
df = pd.DataFrame(columns=['id','Link','Date','Address','Число комнат','Цена','Общая площадь','Жилая площадь кв',
                          'Площадь кухни','Этаж','Этажей в доме', 'Отделка','Расстояние до метро','Год постройки', 'Тип жилья',
                          'Санузел','Балкон','Ремонт'])

k = 0 
for i in tqdm(range(0,to_iter)):
    j += 1
    if j > len(proxies):
        j = 0
#        try:
#            print("Proxy selected: {}".format(proxies[i]))
#            options = webdriver.ChromeOptions()
#            options.add_argument('--proxy-server={}'.format(proxies[i]))
#            driver = webdriver.Chrome(options=options,executable_path='/Users/Nikon/.wdm/drivers/chromedriver/mac64/107.0.5304.62/chromedriver')
#            driver.get("https://www.whatismyip.com/proxy-check/?iref=home")
#            if "Proxy Type" in WebDriverWait(driver, 5).until(EC.visibility_of_element_located((By.CSS_SELECTOR, "p.card-text"))):
#                break
#        except Exception:
#           driver.quit()

    url = 'https://www.cian.ru/sale/flat/'+str(i)
    #options = webdriver.ChromeOptions()
    #options.add_argument('--proxy-server={}'.format(proxies[j]))
    #driver = webdriver.Chrome(options=options,executable_path='/Users/Nikon/.wdm/drivers/chromedriver/mac64/107.0.5304.62/chromedriver')
    driver.get(url)
    time.sleep(5)
    page_source = driver.page_source
    soup = BeautifulSoup(page_source, 'lxml')
    temp_6 = soup.find_all('div', 'a10a3f92e9--container--RXoIe')
    temp = soup.find('h1', 'a10a3f92e9--title--UEAG3')
    #print(soup.find('title').text.split()[0])
    if soup.find('title').text.split()[0] == 'Ошибка':
        continue
    elif temp.text.split()[0].split('-')[0] == 'Участок,':
        continue
    elif soup.find("address").text.strip(',').split()[0] != 'Москва':
        continue
    elif temp_6[0].text == 'Объявление снято с публикации':
        continue
    else:
        df.loc[ len(df.index )] = [np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,np.nan,
                                  np.nan,np.nan,np.nan,np.nan,np.nan]
        df['id'][k] = i
        df['Link'][k] = url
        df['Date'][k] = get_date(soup)
        df['Address'][k] = soup.find("address").text.strip()
        temp = soup.find('h1', 'a10a3f92e9--title--UEAG3')
        df['Число комнат'][k] = temp.text.split()[0].split('-')[0]
        df['Цена'][k] = get_price(soup).replace('\xa0','').replace('₽', '')
        temp = soup.find('h1', 'a10a3f92e9--title--UEAG3')
        df['Общая площадь'][k] = temp.text.split()[2].replace(',', '.')
        df['Расстояние до метро'][k] = find_distance(soup)
        
        if find_living(soup) == 'Не указано':
            df['Жилая площадь кв'][k] = 'Не указано'
        else: 
            df['Жилая площадь кв'][k] = find_living(soup).split()[0].replace(',', '.')
        if find_kitchen(soup) == 'Не указано':
            df['Площадь кухни'][k] = 'Не указано'
        else: 
            df['Площадь кухни'][k] = find_kitchen(soup).split()[0].replace(',', '.')
        if find_floor(soup) == 'Не указано':
            df['Этаж'][k] = 'Не указано'
        else:
            df['Этаж'][k] = (find_floor(soup).split()[0])
                          
        if find_floor(soup) == 'Не указано':
            df['Этажей в доме'][k] = 'Не указано'
        else: 
            df['Этажей в доме'][k] = find_floor(soup).split()[2]
        if find_finishing(soup) == 'Не указано':
            df['Отделка'][k] = 'Не указано'
        else:
            df['Отделка'][k] = find_finishing(soup)
        if find_year(soup) == 'Не указано':
            df['Год постройки'][k] = 'Не указано'
        else:
            df['Год постройки'][k] = find_year(soup)
        if type_floar(soup) == 'Не указано':
            df['Тип жилья'][k] = 'Не указано'
        else:
            df['Тип жилья'][k] = type_floar(soup)
        if type_san(soup) == 'Не указано':
            df['Санузел'][k] = 'Не указано'
        else:
            df['Санузел'][k] = type_san(soup)
        if type_bal(soup) == 'Не указано':
            df['Балкон'][k] = 'Не указано'
        else:
            df['Балкон'][k] = type_bal(soup)
        if type_rem(soup) == 'Не указано':
            df['Ремонт'][k] = 'Не указано'
        else:
            df['Ремонт'][k] = type_rem(soup)
        k += 1

df_all.to_excel('df_all_all.xlsx')
