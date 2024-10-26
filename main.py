import json

f_list = ['first', 'second', 'third', 'four', 'five', 'six', 'seven', 'eight', 'nine', 'ten']


class Dragon:

    def __init__(self, name, param=''):
        self.name = name
        self.param = param

    def printname(self):
        return self.name

    def printinfo(self):
        return "Я обычный дракон, ничего не умею"


class MathDragon(Dragon):

    def matha(self, first, second):
        try:
            self.param = str(first + second)
            return first + second
        except Exception as error:
            return 'Нельзя такое складывать('

    def printinfo(self):
        return "Я дракон, который складывает числа"


class LangDragon(Dragon):

    def lang(self, first, second):
        self.param = str(first + second)
        return str(first) + str(second)

    def printinfo(self):
        return "Я дракон, который складывает строчки"


class DismathDragon(Dragon):

    def dismath(self, first, second):
        try:
            self.param = str(first - second)
            return first - second
        except Exception as err:
            return "Это я вычитать не буду"

    def printinfo(self):
        return "Я дракон, который вычитает числа"


class MultDragon(Dragon):

    def multip(self, first, second):
        self.param = str(first * second)
        return first * second

    def printinfo(self):
        return "Я дракон, который умножает числа"


class DivisDragon(Dragon):

    def division(self, first, second):
        try:
            self.param = str(first / second)
            return first / second
        except Exception as error:
            return "Нельзя так делить!"

    def printinfo(self):
        return "Я дракон, который делит числа"


class StepenDragon(Dragon):

    def stepen(self, first, second):
        try:
            self.param = str(first ** second)
            return first ** second
        except Exception as error:
            return "Нельзя так возводить в степень!"

    def printinfo(self):
        return "Я дракон, который возводит числа в степень"


class BlablaDragon(Dragon):

    def bla(self, c):
        try:
            self.param = str("bla" * c)
            return "bla" * c
        except Exception as err:
            return "Блаблабла с ошибкой!"

    def printinfo(self):
        return "Я делать бла-бла-бла"


class FindLettDragon(Dragon):

    def findletter(self, word, letter):
        self.param = str(word).count(str(letter))
        return str(word).count(str(letter))

    def printinfo(self):
        return "Я дракон, который считает символы в строчке"


class NotlikervowelDragon(Dragon):

    def anti_vowel(self, text):
        try:
            text = list(text)
            for i in text[::-1]:
                if i in 'aeiouAEIOU':
                    text.remove(i)

            self.param = ''.join(text)
            return str(''.join(text))
        except Exception as error:
            return "Так не пойдёт("

    def printinfo(self):
        return "Дркн, ктр ннвдт глсн"


class JsonlikerDragon(Dragon):
    def json_go(self, listik):
        to_json = {'Params': listik}

        with open('Talklog.json', 'w') as f:
            json.dump(to_json, f, indent=4)

    def json_get(self, listik):
        with open('Talklog.json', 'r') as f:
            logs = json.load(f)
            i = 0
            for drparam in listik:
                drparam = logs['Params'][i]
                i += 1



first_dragon = Dragon("Loki")
second_dragon = MathDragon("Mathiania")
third_dragon = LangDragon("Shaman")
four_dragon = DismathDragon("Disoliator")
five_dragon = MultDragon("Multik")
six_dragon = DivisDragon("Komanda")
seven_dragon = StepenDragon("Stepan")
eight_dragon = BlablaDragon("Car")
nine_dragon = FindLettDragon("Yandex")
ten_dragon = NotlikervowelDragon("Sglsn")

teh_dragon = JsonlikerDragon("Jsonya")

first_dragon.printinfo()

second_dragon.printinfo()
print(second_dragon.matha(2, 3))

third_dragon.printinfo()
print(third_dragon.lang(2, 3))

four_dragon.printinfo()
print(four_dragon.dismath(10, 5))

five_dragon.printinfo()
print(five_dragon.multip(2, 4))

six_dragon.printinfo()
print(six_dragon.division(10, 2))

seven_dragon.printinfo()
print(seven_dragon.stepen(2, 5))

eight_dragon.printinfo()
print(eight_dragon.bla(3))

nine_dragon.printinfo()
print(nine_dragon.findletter('aaa ih tyt tri', 'a'))

ten_dragon.printinfo()
print(ten_dragon.anti_vowel('tyt neskolko glasnih'))

s_list = [first_dragon.param, second_dragon.param, third_dragon.param, four_dragon.param, five_dragon.param,
          six_dragon.param, seven_dragon.param, eight_dragon.param, nine_dragon.param, ten_dragon.param]

teh_dragon.json_go(s_list)

print(five_dragon.multip(100, 100))

print(five_dragon.param)


hst = teh_dragon.json_get(s_list)





