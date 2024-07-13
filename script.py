import json
import os
import io
import re
import sys
import zlib
import time
import base64
import shutil
import hashlib
from pathlib import Path
from datetime import datetime

import argparse
import colorama
from colorama import Fore

#все ради винды
colorama.init(autoreset=True)
if sys.platform == "win32":
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

############################################################

SAMPLE_INFO = {
    "iso": "",  # ISO 639-2
    "date": "",  # ISO 8601
    "hash": "",  # BLAKE2b
    "img": ""  # data:image
}
SAMPLE_FILE = {
    "info": {},
    "data": {}
}
SAMPLE_TEXT = "This file is generated automatically and any changes made to it may be overwritten. Please do not modify the contents of this file manually."

START_TIME = None
DEF_FILE_ROUT_NAME = "languages"

IS_SAFE = False

def hashBlake2b(data, digest_size=16):
    '''Створюємо об'єкт хешу з алгоритмом BLAKE2b'''

    if isinstance(data, str):
        data = data.encode('utf-8')

    hasher = hashlib.blake2b(digest_size=digest_size)
    hasher.update(data)
    hash_hex = hasher.hexdigest()

    return hash_hex

############################################################

#   Класс для создания файла интеграции языковых деревьев
class ReactTS:
    _infoObj = "interface LangInfoObj {date: string; hash: string; img: string; name: { short: string; english: string; native: string; }}"
    _langObj = "interface LangObj {info: LangInfoObj; data: LangDataObj;}"
    _dataObjName = "LangDataObj"

    isJS=False

    _dir = ""

    _fileMap = {}
    _dataStruct = ""

    def __init__(self, dirPath=""):
        self._dir = dirPath

    def addObj(self, jsonStringData=""):
        '''добавление одного из языковых деревьев для генерации структуры с него'''
        self._dataStruct = jsonStringData

    def add(self, code="", fileName=""):
        '''Добавление языкового файла'''
        self._fileMap[code] = fileName

    def createFile(self):
        '''Создание файла интеграции'''
        if not self.isJS:
            filePath = Path(self._dir) / (DEF_FILE_ROUT_NAME+".ts")
        else:
            filePath = Path(self._dir) / (DEF_FILE_ROUT_NAME+".js")

        if IS_SAFE:
            if os.path.exists(filePath):
                echo.info(echo.pink("SafeMode"))
                shutil.copy(filePath, str(filePath)+".old")

        if not self.isJS:
            echo.info("Create React-integration with TypeScript file "+echo.orange(DEF_FILE_ROUT_NAME+".ts"))
        else:
            echo.info("Create integration with JavaScript file "+echo.orange(DEF_FILE_ROUT_NAME+".js"))

        with open(filePath, 'w', encoding="utf-8") as file:
            file.write("/* "+SAMPLE_TEXT+" */\n")#коментарий что сгенерировано автоматически

            #добавляем импорт
            for code, langFileName in self._fileMap.items():
                file.write("import fileLang"+code.upper()+" from './"+langFileName+"'\n")
            file.write("\n")

            #интегрируем языковое дерево (структура)
            if not self.isJS:
                file.write(self._infoObj+"\n")
                file.write("interface "+self._dataObjName+" { "+self._generateStruct()+" }\n")
                file.write(self._langObj+"\n")
                file.write("interface LanguagesTreeObj {")
                for code, _ in self._fileMap.items():
                    file.write(code+": LangObj; ")
                file.write("}\n")

            file.write("\n")

            #Формирем глобальную структуру языкового дерева
            if not self.isJS:
                file.write("export const Languages: LanguagesTreeObj = {\n")
            else:
                file.write("export var LanguagesTreeObj = {\n")
            for code, imgData in self._fileMap.items():
                file.write("\t"+code+": {\n")
                if True:
                    file.write("\t\tinfo: {\n")
                    if True:
                        file.write("\t\t\timg: fileLang"+code.upper()+".info.img,\n")
                        file.write("\t\t\tdate: fileLang"+code.upper()+".info.date,\n")
                        file.write("\t\t\thash: fileLang"+code.upper()+".info.hash,\n")
                        file.write("\t\t\tname: {\n")
                        if True:
                            file.write("\t\t\t\tshort: \""+code.upper()+"\",\n")
                            file.write("\t\t\t\tenglish: \""+LangCodes.ARR[code]["english"]+"\",\n")
                            file.write("\t\t\t\tnative: \""+LangCodes.ARR[code]["native"]+"\",\n")
                        file.write("\t\t\t},\n")
                    file.write("\t\t},\n")
                    file.write("\t\tdata: fileLang"+code.upper()+".data,\n")
                file.write("\t},\n")
            file.write("\n}\n")
            file.close()

    def _generateStruct(self):
        ''' Генерация структуры обьекта для TS'''

        def recursive(data):
            if isinstance(data, dict):# Если это словарь, рекурсивно обрабатываем каждое значение
                result = []

                for key, value in data.items():
                    if isinstance(value, dict):
                        result.append(key+": {"+recursive(value)+"},")

                    elif isinstance(value, (dict, list)):
                        result.append(key+": string[];")

                    else:
                        result.append(key+": string;")

                return ' '.join(map(str, result))


        data = json.loads(self._dataStruct)
        struct = recursive(data)
        return struct



#   Класс вывода в консоль
class echo:
    _PINK = Fore.MAGENTA
    _OKBLUE = Fore.BLUE
    _OKCYAN = Fore.CYAN
    _OKGREEN = Fore.GREEN
    _ORANGE = Fore.YELLOW
    _RED = Fore.RED
    _BOLD = '\033[1m'
    _UNDERLINE = '\033[4m'
    _ENDC = '\033[0m'
    _SLOT = ""

    @staticmethod
    def format(text="", bold=False, underline=False):
        if bold:
            text = echo._BOLD+text+echo._ENDC
        if bold:
            text = echo._UNDERLINE+text+echo._ENDC
        return text

    @staticmethod
    def blue(text="", bold=False, underline=False):
        text = echo._OKBLUE+text+echo._ENDC
        return echo.format(text, bold, underline)

    @staticmethod
    def green(text="", bold=False, underline=False):
        text = echo._OKGREEN+text+echo._ENDC
        return echo.format(text, bold, underline)

    @staticmethod
    def cyan(text="", bold=False, underline=False):
        text = echo._OKCYAN+text+echo._ENDC
        return echo.format(text, bold, underline)

    @staticmethod
    def pink(text="", bold=False, underline=False):
        text = echo._PINK+text+echo._ENDC
        return echo.format(text, bold, underline)

    @staticmethod
    def orange(text="", bold=False, underline=False):
        text = echo._ORANGE+text+echo._ENDC
        return echo.format(text, bold, underline)

    @staticmethod
    def red(text="", bold=False, underline=False):
        text = echo._RED+text+echo._ENDC
        return echo.format(text, bold, underline)
    #

    @staticmethod
    def groupStart(text=""):
        min, sec, millsec = echo.timePoint()
        print("["+f"{min:02d}:{sec:02d}.{millsec:04d}"+"] "+echo._SLOT + "▶ " + echo._OKCYAN + echo._BOLD+text+ echo._ENDC)
        echo._SLOT = echo._SLOT+">"

    @staticmethod
    def groupEnd():
        echo._SLOT = echo._SLOT[:-1] if echo._SLOT else echo._SLOT


    #

    @staticmethod
    def timePoint():
        '''Время выполнения программы с момента начала'''

        end = time.time()
        elapsed = end - START_TIME

        minutes = int(elapsed // 60)
        seconds = int(elapsed % 60)
        milliseconds = int((elapsed - int(elapsed)) * 1000)

        return minutes, seconds, milliseconds

    @staticmethod
    def info(info=""):
        '''Вывод текста в консоль с отступом'''
        min, sec, millsec = echo.timePoint()
        print("["+f"{min:02d}:{sec:02d}.{millsec:04d}"+"] "+echo._SLOT + " " + info)

    @staticmethod
    def error(info="", msg=""):
        '''Вывод ошибки в консоль'''
        echo.info(echo._RED + echo._BOLD + info + echo._ENDC + "\t" + echo._ORANGE + msg + echo._ENDC)

    @staticmethod
    def panic(info="", msg=""):
        '''Вывод ошибки и завершение работы'''
        echo.error(info, msg)
        sys.exit(1)

#   Класс проверки ISO-кодов языка
class LangCodes:
    CONST = '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'
    ARR = {}

    @staticmethod
    def isset(code=""):
        '''Проверка на существование такой локали'''

        if code in LangCodes.ARR:
            return True
        else:
            return False

    @staticmethod
    def init():
        '''Инициализация буферезированного дерево локалей (ISO 639-2)'''
        LangCodes.ARR = LangCodes.decode(LangCodes.CONST)

    @staticmethod
    def decode(base64string=""):
        b = base64.b64decode(base64string)
        d = zlib.decompress(b)
        j = d.decode('utf-8')
        return json.loads(j)

    @staticmethod
    def encode(arr={}):
        j = json.dumps(arr, ensure_ascii=False)
        b = j.encode('utf-8')
        d = zlib.compress(b, level=9)
        return base64.b64encode(d)

# класс работы с файлами
class FileSystem:

    @staticmethod
    def validFilePath(filePath=""):
        '''проверка пути к файлу'''

        if not os.path.isfile(filePath):
            echo.error("This file does not exist", filePath)
            return False

        _, extension = os.path.splitext(filePath)
        if not extension.lower() == ".json":
            echo.error("Unsupported type, only json is allowed", filePath)
            return False

        return True

    @staticmethod
    def validDirPath(dirPath=""):
        '''проверка пути к директории'''

        if not os.path.isdir(dirPath):
            echo.error("This dir does not exist", dirPath)
            return False

        return True


    #

    @staticmethod
    def parseJSON(filePath=""):
        '''Парсим JSON'''

        with open(filePath, 'r', encoding="utf-8") as file:
            return json.load(file)

    @staticmethod
    def parseDir(dir=""):
        '''Получение имен файлов языкового дерева'''
        items = os.listdir(dir)

        files = [item for item in items
                 if item.endswith('.json')
                 and item != DEF_FILE_ROUT_NAME+".json"
                 and os.path.isfile(os.path.join(dir, item))]

        return files

    @staticmethod
    def dirFromFile(filePath=""):
        '''Получение пути к директории по пути к файлу'''
        return Path(filePath).parent

    @staticmethod
    def fileFromFile(filePath=""):
        '''Получение названия файла по пути к файлу'''
        return os.path.basename(filePath)

#Класс работы с языковым деревом
class Tree:

    @staticmethod
    def generateFileRouts(dirPath="/", langs={}):
        '''Создание файла роутов на языки '''
        fileName = DEF_FILE_ROUT_NAME+".json"

        echo.groupStart("Create "+fileName)

        fileLangBuf = {}
        filePath = Path(dirPath) / fileName

        for fileName, langJSON in langs.items():
            lang = json.loads(langJSON)

            lang = lang["info"]
            buf = LangCodes.ARR[lang["iso"]]

            echo.info("Add: "+echo.pink(lang["iso"])+" "+echo.orange(buf["english"]))
            fileLangBuf[lang["iso"]] = {
                "file": fileName,
                "img": lang["img"],
                "name": {
                    "short": lang["iso"].upper(),
                    "full": buf["english"],
                    "native": buf["native"]
                }
            }

        if IS_SAFE:
            if os.path.exists(filePath):
                echo.info(echo.pink("SafeMode"))
                shutil.copy(filePath, str(filePath)+".old")

        echo.groupEnd()
        with open(filePath, 'w', encoding="utf-8") as fileTY:
            json.dump(fileLangBuf, fileTY, indent=2, ensure_ascii=False)
            fileTY.close()

    @staticmethod
    def generateFile(filePath="/", langObj={}):
        '''Создание языкового файла '''

        if not "info" in langObj or not "data" in langObj:
            echo.error("This is not a language tree object", str({"info": "info" in langObj, "data": "data" in langObj}))
            return

        #Генерация хеша по наполнению
        hashString = json.dumps(langObj["data"], ensure_ascii=False)
        hashString = hashBlake2b(hashString)

        if hashString == langObj["info"]["hash"]:#Информирование если файл не поменялся
            echo.info(echo.orange("File hasn't changed"))

        else:
            langObj["info"]["hash"] = hashString
            langObj["info"]["date"] = datetime.now().date().isoformat()
            echo.info(echo.green("The file has been recompiled"))

            if IS_SAFE:
                if os.path.exists(filePath):
                    echo.info(echo.pink("SafeMode"))
                    shutil.copy(filePath, str(filePath)+"."+langObj["info"]["hash"]+".old")

        with open(filePath, 'w', encoding="utf-8") as fileGH:
            json.dump(langObj, fileGH, indent=2, ensure_ascii=False)
            fileGH.close()

    @staticmethod
    def parseAndValidFile(filePath = ""):
        '''Чтение данных из файла, проверки и формирование валидной структуры'''
        fileObj = SAMPLE_FILE

        echo.groupStart("Read and Valid: "+echo.blue(filePath))
        file = FileSystem.parseJSON(filePath)

        ####

        if not "info" in file: #Отсекаем если нет INFO
            echo.error("Block `info` not found")
            echo.groupEnd()
            return False, fileObj
        if not "data" in file: #Отсекаем если нет тела
            echo.error("Block `data` not found")
            echo.groupEnd()
            return False, fileObj

        status, info = Tree.validInfoBlock(file["info"])
        if not status: #Отсекаем если INFO невалидный
            echo.groupEnd()
            return False, fileObj

        ####

        data = Tree.validDataBlock(file["data"])
        if not Tree.validDataKeys(data): #Отсекаем если есть невалидные ключи
            echo.groupEnd()
            return False, fileObj

        data = Tree.sortDataBlock(data) #Сортируем по ключам

        fileObj["info"] = info
        fileObj["data"] = data
        echo.info(echo.green("File valid"))

        echo.groupEnd()
        return True, fileObj

    ##

    @staticmethod
    def validInfoBlock(fileInfo={}):
        '''Валидация и парсинг информационного блока'''
        info = SAMPLE_INFO

        if not "iso" in fileInfo:
            echo.error("Field `info.iso` not found")
            return False, info
        if not LangCodes.isset(fileInfo["iso"]):
            echo.error("Field `info.iso` is invalid (ISO 639-2)", fileInfo["iso"])
            return False, info
        info["iso"] = fileInfo["iso"]

        if "date" in fileInfo:
            info["date"] = fileInfo["date"]
        if "hash" in fileInfo:
            info["hash"] = fileInfo["hash"]

        if "img" in fileInfo:
            pattern = re.compile(r'^data:image/[a-zA-Z0-9.+-]+;base64,[A-Za-z0-9+/]+={0,2}$')
            if bool(pattern.match(fileInfo["img"])) or fileInfo["img"] == "" : #Проверка по соответсвию на data:image
                info["img"] = fileInfo["img"]
            else:
                echo.error("Field `info.img` is invalid", fileInfo["img"])
                return False, info

        return True, info

    @staticmethod
    def validDataBlock(data):
        '''Проверка данных и удаление лишнего'''
        data = Tree._validDataBlock(data)
        data = Tree._validDataClearArrays(data)
        return data


    @staticmethod
    def _validDataBlock(data):
        '''Рекурсивная проверка языкового дерева и отсечение всего лишнего (все кроме строк, обьектов и массивов(не пустых) будет отброшено)'''

        if isinstance(data, dict):# Если это словарь, рекурсивно обрабатываем каждое значение
            result = {}

            for key, value in data.items():
                Tree.bufAdr.append(key)#Инкремент позиции

                if isinstance(value, (dict, list)):
                    if len(value) == 0:
                        echo.error("Variable is empty", Tree._adrToString())
                    else:
                        result[key] = Tree._validDataBlock(value)

                elif isinstance(value, str):
                    result[key] = Tree._validDataBlock(value)

                else:
                    echo.error("Invalid type value", Tree._adrToString())

                Tree.bufAdr.pop()#декремент позиции

            return result

        elif isinstance(data, list):# Если это список, рекурсивно обрабатываем каждый элемент
            result = []

            for item in data:
                if isinstance(item, str) or isinstance(item, (dict, list)):
                    result.append(Tree._validDataBlock(item))
                else:
                    echo.error("Invalid type value in array", Tree._adrToString("["+str(item)+"]"))

            return result

        elif isinstance(data, str):# Если это строка, возвращаем ее
            return data

        # Все остальные типы данных игнорируем
        echo.error("Invalid type value", str(data))
        return None


    @staticmethod
    def validDataKeys(data):
        return Tree._validDataKeys(data, True)
    
    @staticmethod
    def _validDataKeys(data, status):
        '''Рекурсивная проверка на соответствие ключей'''
        regex = r'^[a-zA-Z_][a-zA-Z0-9_]*$'

        if isinstance(data, dict):
            for key, value in data.items():
                Tree.bufAdr.append(key)#Инкремент позиции

                if not re.match(regex, key):
                    echo.error("Invalid key "+echo.pink(key), Tree._adrToString())
                    status = False

                Tree._validDataKeys(value, status)
                Tree.bufAdr.pop()#декремент позиции

        elif isinstance(data, list):
            for item in data:
                Tree._validDataKeys(item, status)

        return status


    @staticmethod
    def _validDataClearArrays(data):
        '''Удаление лишних елементов из массивов'''

        for key, value in list(data.items()):
            Tree.bufAdr.append(key)#Инкремент позиции

            if isinstance(value, dict):
                data[key] = Tree._validDataClearArrays(value)

            elif isinstance(value, list):
                original_list = value[:]

                data = {}
                data[key] = []
                for item in value:
                    if isinstance(item, str):
                        data[key].append(item)

                index = 0
                for item in original_list:
                    if item not in data[key]:
                        echo.error("Incorrect data type in array", Tree._adrToString(str(index)))
                        index += 1

            elif isinstance(value, str):
                data[key] = value

            else:
                echo.error("Incorrect data type", Tree._adrToString())

            Tree.bufAdr.pop()#декремент позиции
        return data

    @staticmethod
    def sortDataBlock(data):
        '''сортировка по ключам языкового дерева'''

        if isinstance(data, dict):
            sorted_data = {}

            for key, v in sorted(data.items()):
                Tree.bufAdr.append(key)#Инкремент позиции

                if v is not None:  # Перевірка, чи значення не є None
                    bufValue = Tree.sortDataBlock(v)
                    if isinstance(v, str) or bool(bufValue):
                        sorted_data[key] = bufValue

                else:
                    echo.error("Variable is empty", Tree._adrToString())

                Tree.bufAdr.pop()#декремент позиции

            return sorted_data


        elif isinstance(data, list):
            sorted_list = []
            pos = 0

            for item in data:
                if item is not None:  # Перевірка, чи елемент не є None
                    sorted_item = Tree.sortDataBlock(item)
                    sorted_list.append(sorted_item)

                else:
                    echo.error("Variable is empty", Tree._adrToString("["+str(pos)+"]"))
                pos += 1

            return sorted_list

        return data
    
    ####

    bufAdr = []
    _bufClone = {
        "file": "",
        "data": {}
    }
    @staticmethod
    def _adrToString(key=""):
        '''Формирование строки "адреса" переменной'''
        bufArr = Tree.bufAdr

        if key != "":
            if not key[0] == "[":
                key = "."+key

        return ".".join(map(str, bufArr))+key

    @staticmethod
    def _compareAndClone(obj, obj2):
        '''Сравнение двух обьектов и получение клона по эталону'''

        if isinstance(obj, dict):#проверка на обьект
            if not isinstance(obj2, dict):
                echo.error("Data type mismatch", "[dict]"+Tree._adrToString())#несопали типы. ожидался обьект
                return obj

            ret = {}
            for key, value in obj.items():
                Tree.bufAdr.append(key)#Инкремент позиции

                if key in obj2:
                    ret[key] = Tree._compareAndClone(value, obj2[key])
                else:
                    echo.error("++++ Variable is missing", Tree._adrToString())#Нехватает переменной
                    ret[key] = value

                Tree.bufAdr.pop()#декремент позиции

            for key, value in obj2.items():
                if not key in ret:
                    echo.error("---- Extra variable", Tree._adrToString(key))#такой переменной нет в эталоне
            return ret

        elif isinstance(obj, list):#Проверка на массив
            if not isinstance(obj, list):
                echo.error("Data type mismatch", "[list]"+Tree._adrToString())#Несовпали типы. ожидался массив
                return obj
            else:
                return obj2

        elif isinstance(obj, str):# проверка на строку
            if not isinstance(obj2, str):
                echo.error("Data type mismatch", "[str]"+Tree._adrToString())#Несовпали типы. ожидалась строка
                return obj
            else:
                return obj2

        else:
            echo.error("Break recursive", "[Invalid type]"+Tree._adrToString())#Обший критикал

    @staticmethod
    def _compareAndCloneFiles(filePath="", patternFilePath=""):
        '''Сравнение файлов и получение копии с зашитой от пиздеца питона'''

        if not Tree._bufClone["file"] == FileSystem.fileFromFile(patternFilePath):#Проверка на наличие буферизированых данных
            status, obj = Tree.parseAndValidFile(patternFilePath)

            if not status:
                echo.error("Break compare")
                return False, {}

            Tree._bufClone["file"] = FileSystem.fileFromFile(patternFilePath)
            Tree._bufClone["data"] = obj["data"]
            obj = None

        status, obj = Tree.parseAndValidFile(filePath)
        if not status:
            echo.error("Break compare")
            return False, {}

        echo.groupStart("Compare "+echo.pink(Tree._bufClone["file"])+echo.cyan(" >> ")+echo.blue(FileSystem.fileFromFile(filePath)))
        obj["data"] = Tree._compareAndClone(Tree._bufClone["data"], obj["data"])
        echo.groupEnd()
        return True, obj

#класс обработки консольных вызовов
class ConsoleRouts:

    @staticmethod
    def fileOrDir(filePath, dirPath, fileFunc, dirFunc, **kwargs):
        if not filePath == None:
            filePath = filePath.strip().strip('"').strip("'")
            echo.groupStart("File mode ["+echo.blue(str(filePath))+"]")

            if FileSystem.validFilePath(filePath):
                fileFunc(filePath, **kwargs)

            echo.groupEnd()

        elif not dirPath == None:
            dirPath = dirPath.strip().strip('"').strip("'")
            echo.groupStart("Dir mode ["+echo.blue(str(dirPath))+"]")

            if FileSystem.validDirPath(dirPath):
                dirFunc(dirPath, **kwargs)

            echo.groupEnd()

        else:
            echo.error("Error! You must use one of the parameters:", "[`--file` | `--dir` ]")

    ########

    @staticmethod
    def _codeOnly(obj={}, codeOnly=False):
        '''вывод только массива кодов или готовой структуры'''
        if not codeOnly:
            echo.info(json.dumps(obj, ensure_ascii=False))

        else:
            buf = []
            for code, _ in obj.items():
                buf.append(code)
            echo.info(json.dumps(buf, ensure_ascii=False))


    @staticmethod
    def searchFromCodes(code="", codeOnly=False):
        '''поиск по коду'''
        code = code.strip().strip('"').strip("'")

        result = {}
        for key, obj in LangCodes.ARR.items():
            if code in key:
                result[key] = obj

        ConsoleRouts._codeOnly(result, codeOnly)

    @staticmethod
    def searchFromNames(name="", codeOnly=False):
        '''поиск по названию'''
        name = name.strip().strip('"').strip("'")
        name = name.lower()

        result = {}
        for key, obj in LangCodes.ARR.items():
            if name in obj["english"].lower() or name in obj["native"].lower():
                result[key] = obj

        ConsoleRouts._codeOnly(result, codeOnly)

    @staticmethod
    def fileCompress(filepath=""):
        '''чтение json по ссылке и возрат сжатой строки'''
        filepath = filepath.strip().strip('"').strip("'")

        if  FileSystem.validFilePath(filepath):
            data = FileSystem.parseJSON(filepath)
            base = LangCodes.encode(data)
            echo.info(str(base))

    @staticmethod
    def printCodes(codeOnly=False):
        '''получить массив кодов'''
        ConsoleRouts._codeOnly(LangCodes.ARR, codeOnly)


    #####

    @staticmethod
    def _checkFileFromFile(filePath):
        '''проверка конкретного файла'''
        Tree.parseAndValidFile(filePath)

    @staticmethod
    def _checkFileFromDir(dirPath):
        '''проверка всех файлов в директории'''
        filesInDir = FileSystem.parseDir(dirPath)
        dirBuf = Path(dirPath)

        for fileName in filesInDir:
            fullPath = str(dirBuf / fileName)

            if FileSystem.validFilePath(fullPath):
                Tree.parseAndValidFile(fullPath)

    @staticmethod
    def checkFile(filePath, dirPath):
        '''проверяем файл(ы) на валидность'''
        ConsoleRouts.fileOrDir(
            filePath=filePath,
            dirPath=dirPath,
            fileFunc=ConsoleRouts._checkFileFromFile,
            dirFunc=ConsoleRouts._checkFileFromDir
        )

    #####

    @staticmethod
    def _cloneFileFromFile(filePath, pattern):
        '''клонирование между двух файлов'''
        status, bufFile = Tree._compareAndCloneFiles(str(filePath), str(pattern))
        if status:
            Tree.generateFile(filePath, bufFile)
        else:
            echo.error("File not clone")

    @staticmethod
    def _cloneFileFromDir(dirPath, pattern):
        '''клонирование ко всей директории'''

        filesInDir = FileSystem.parseDir(dirPath)
        dirBuf = Path(dirPath)

        for fileName in filesInDir:
            fullPath = str(dirBuf / fileName)
            if not fullPath == pattern:

                status, bufFile = Tree._compareAndCloneFiles(fullPath, str(pattern))
                if status:
                    Tree.generateFile(fullPath, bufFile)
                else:
                    echo.error("File not clone")


    @staticmethod
    def cloneFiles(filePath, dirPath, pattern):
        '''клонируем файл(ы) по еталону'''
        ConsoleRouts.fileOrDir(
            filePath=filePath,
            dirPath=dirPath,
            fileFunc=ConsoleRouts._cloneFileFromFile,
            dirFunc=ConsoleRouts._cloneFileFromDir,
            pattern=pattern
        )

    #####

    @staticmethod
    def _compareFileFromFile(filePath, pattern):
        '''сравнение с конкретным файлом'''
        Tree._compareAndCloneFiles(str(filePath), str(pattern))


    @staticmethod
    def _compareFileFromDir(dirPath, pattern):
        '''сравнить со всеми файлами в директории'''

        filesInDir = FileSystem.parseDir(dirPath)
        dirBuf = Path(dirPath)

        for fileName in filesInDir:
            fullPath = str(dirBuf / fileName)
            if not fullPath == pattern:
                Tree._compareAndCloneFiles(fullPath, str(pattern))

    @staticmethod
    def compareFiles(filePath, dirPath, pattern):
        '''Сверяем файл(ы) по еталону'''
        ConsoleRouts.fileOrDir(
            filePath=filePath,
            dirPath=dirPath,
            fileFunc=ConsoleRouts._compareFileFromFile,
            dirFunc=ConsoleRouts._compareFileFromDir,
            pattern=pattern
        )

    #

    @staticmethod
    def _routsFileFromFile(filePath):
        '''создаем файл роута по конкретному файлу'''
        langs = {}

        status, treeObj = Tree.parseAndValidFile(filePath)
        if status:
            dirPath = FileSystem.dirFromFile(filePath)
            fileName = FileSystem.fileFromFile(filePath)

            langs[fileName] = treeObj
            Tree.generateFileRouts(str(dirPath), langs)

    @staticmethod
    def _routsFileFromDir(dirPath):
        '''создаем файл роута по всем файлам в директории'''

        filesInDir = FileSystem.parseDir(dirPath)
        dirBuf = Path(dirPath)
        langs = {}
        bufISO=[]

        for fileName in filesInDir:
            fullPath = str(dirBuf / fileName)
            status = False
            fuckingLangTree = None

            status, fuckingLangTree = Tree.parseAndValidFile(fullPath)
            if status:

                if fuckingLangTree["info"]["iso"] in bufISO: #отсекаем если уже есть такой язык
                    echo.error("Duplicate", fuckingLangTree["info"]["iso"])
                    continue

                bufISO.append(fuckingLangTree["info"]["iso"])
                langs[fileName] = json.dumps(fuckingLangTree, ensure_ascii=False)

        Tree.generateFileRouts(str(dirPath), langs)

    @staticmethod
    def routsFile(filePath, dirPath):
        '''создаем файл роута'''
        ConsoleRouts.fileOrDir(
            filePath=filePath,
            dirPath=dirPath,
            fileFunc=ConsoleRouts._routsFileFromFile,
            dirFunc=ConsoleRouts._routsFileFromDir
        )

    #

    @staticmethod
    def _treeFileFromFile(filePath):
        '''перекомпиляция конкретного файла (все проблемые места автоматически удаляются)'''
        _, treeObj = Tree.parseAndValidFile(filePath)
        Tree.generateFile(filePath, treeObj)


    @staticmethod
    def _treeFileFromDir(dirPath):
        '''Перекомпиляция всех файлов в директории (все проблемые места автоматически удаляются)'''
        filesInDir = FileSystem.parseDir(dirPath)
        dirBuf = Path(dirPath)

        for fileName in filesInDir:
            fullPath = str(dirBuf / fileName)

            if FileSystem.validFilePath(fullPath):
                _, treeObj = Tree.parseAndValidFile(fullPath)
                Tree.generateFile(fullPath, treeObj)


    @staticmethod
    def treeFile(filePath, dirPath):
        '''перекомпилируем файл(ы) языкового дерева'''
        ConsoleRouts.fileOrDir(
            filePath=filePath,
            dirPath=dirPath,
            fileFunc=ConsoleRouts._treeFileFromFile,
            dirFunc=ConsoleRouts._treeFileFromDir
        )

    ####

    @staticmethod
    def _reactTypeScriptFromFile(filePath):
        '''перекомпиляция конкретного файла для интеграции с React-TS'''
        _, treeObj = Tree.parseAndValidFile(filePath)
        react = ReactTS(str(FileSystem.dirFromFile(filePath)))
        fileName = FileSystem.fileFromFile(filePath)

        react.add(
            code=treeObj["info"]["iso"],
            fileName=fileName,
        )
        react.addObj(jsonStringData=json.dumps(treeObj["data"], ensure_ascii=False))

        react.createFile()


    @staticmethod
    def _reactTypeScriptFromDir(dirPath):
        '''Перекомпиляция всех файлов в директории для в интеграции с React-TS'''
        filesInDir = FileSystem.parseDir(dirPath)
        dirBuf = Path(dirPath)
        react = ReactTS(str(dirBuf))
        bufISO=[]

        for fileName in filesInDir:
            fullPath = str(dirBuf / fileName)

            _, treeObj = Tree.parseAndValidFile(fullPath)
            if treeObj["info"]["iso"] in bufISO: #отсекаем если уже есть такой язык
                echo.error("Duplicate", treeObj["info"]["iso"])
                continue

            bufISO.append(treeObj["info"]["iso"])
            react.add(
                code=treeObj["info"]["iso"],
                fileName=fileName,
            )
            react.addObj(jsonStringData=json.dumps(treeObj["data"], ensure_ascii=False))

        react.createFile()

    @staticmethod
    def reactTypeScriptFile(filePath, dirPath):
        '''Компилируем файл(ы) языкового дерева в интеграцию для React-TS'''
        ConsoleRouts.fileOrDir(
            filePath=filePath,
            dirPath=dirPath,
            fileFunc=ConsoleRouts._reactTypeScriptFromFile,
            dirFunc=ConsoleRouts._reactTypeScriptFromDir
        )

###########################################################

if __name__ == '__main__':
    START_TIME = time.time()
    consoleParser = argparse.ArgumentParser(description="Script for working with language trees. Checking, comparing, recompiling, cloning changes, etc.")

    consoleParser.add_argument('--file', '-f', type=str, required=False, help="Path to file")
    consoleParser.add_argument('--dir', '-d', type=str, required=False, help="Path to directory")
    consoleParser.add_argument('--safe', '-s', action='store_true', required=False, help="Safe mode. When overwritten, a copy of the old version will be created.")

    check = consoleParser.add_argument_group('Checking the correctness of files and displaying all errors. Mandatory use [`--file` OR `--dir` ]')
    check.add_argument('--check', '-ch', action='store_true', help="Language file checking mode")

    compiling = consoleParser.add_argument_group('Compilation of language files according to uniform rules. Mandatory use [`--file` OR `--dir` ]')
    compiling.add_argument('--compiling', '-cp', action='store_true', help="File compilation mode")
    compiling.add_argument('--tree', action='store_true', help="Recompiling the language file(s)")
    compiling.add_argument('--routs', action='store_true', help="Generating a route file")
    compiling.add_argument('--reactTS', action='store_true', help="Compiling a file for React-integration with TypeScript")

    clone = consoleParser.add_argument_group('Working with cloning. Compare files and apply changes relative to the `reference`. Mandatory use [`--file` OR `--dir` ]  [`--clone` OR `--compare` ]')
    clone.add_argument('--clone', '-cl', action='store_true', help="Cloning mode")
    clone.add_argument('--compare', action='store_true', help="Comparing a reference with a file or files in a directory")
    clone.add_argument('--pattern', type=str, help="Reference file")

    codes = consoleParser.add_argument_group('Obtaining information on embedded languages. Without [`--code` OR `--name` OR `--fileCompress` ] displays the entire list of embedded codes')
    codes.add_argument('--iso', action='store_true', help="Mode of operation with embedded languages")
    codes.add_argument('--codeOnly', action='store_true', help="Show only codes")
    codes.add_argument('--code', type=str, help="Поиск по кодам")
    codes.add_argument('--name', type=str, help="Search by name")
    codes.add_argument('--fileJSON', type=str, help="Get base64 compressed from json file")

    consoleARG = consoleParser.parse_args()
    IS_SAFE = consoleARG.safe
    LangCodes.init()

    if consoleARG.check:
        echo.groupStart("Language file checking mode")
        ConsoleRouts.checkFile(filePath=consoleARG.file, dirPath=consoleARG.dir)
        echo.groupEnd()


    elif consoleARG.compiling:
        echo.groupStart("File compilation mode")

        if consoleARG.routs:#создаем файл роута
            ConsoleRouts.routsFile(filePath=consoleARG.file, dirPath=consoleARG.dir)

        elif consoleARG.tree:#перекомпилируем файл(ы) языкового дерева
            ConsoleRouts.treeFile(filePath=consoleARG.file, dirPath=consoleARG.dir)

        elif consoleARG.reactTS:#перекомпилируем файл(ы) для тайпскрипта
            ConsoleRouts.reactTypeScriptFile(filePath=consoleARG.file, dirPath=consoleARG.dir)

        else:
            echo.error("Error! You must use one of the parameters", "[ --tree | --routs ]")
        echo.groupEnd()


    elif consoleARG.clone or consoleARG.compare:
        if consoleARG.pattern == None:
            echo.error("Error! You need to use the parameter", "[`--pattern`]")

        else:
            consoleARG.pattern = consoleARG.pattern.strip().strip('"').strip("'")
            echo.info("File pattern ["+echo.pink(str(consoleARG.pattern), True)+"]")

            if consoleARG.clone:#клонирование изменений из одного файла
                echo.groupStart("Cloning mode")
                ConsoleRouts.cloneFiles(filePath=consoleARG.file, dirPath=consoleARG.dir, pattern=consoleARG.pattern)

            elif consoleARG.compare:#Сверяем файл(ы) по еталону'
                echo.groupStart("Compare mode")
                ConsoleRouts.compareFiles(filePath=consoleARG.file, dirPath=consoleARG.dir, pattern=consoleARG.pattern)

            else:
                echo.error("Error! You must use one of the parameters", "[ --clone | --compare ]")
            echo.groupEnd()


    elif consoleARG.iso:
        echo.groupStart("Mode of operation with embedded languages")

        if not consoleARG.code == None:#поиск по коду
            ConsoleRouts.searchFromCodes(consoleARG.code, consoleARG.codeOnly)

        elif not consoleARG.name == None:#поиск по названию
            ConsoleRouts.searchFromNames(consoleARG.name, consoleARG.codeOnly)

        elif not consoleARG.fileJSON == None:#чтение json по ссылке и возрат сжатой строки
            ConsoleRouts.fileCompress(consoleARG.fileJSON)

        else:#получить массив кодов
            ConsoleRouts.printCodes(consoleARG.codeOnly)
        echo.groupEnd()

    else:
        echo.error("Error! You must use one of the parameters", "[ --check | --compiling | --clone | --compare | --iso ]")
