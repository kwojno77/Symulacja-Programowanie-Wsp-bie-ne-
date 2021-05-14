# Symulacja-Programowanie-Wspóbieżne
Program symuluje przepływ pakietów pomiędzy wierzchołkami grafu skierowanego. Wierzchołki połączone są kolejno, tzn. wierzchołek 1 wysyła pakiet do wierzchołka 2, ten z kolei do wierzchołka 3 itd. W programie osobne wątki: nadawcę, wierzchołki pośrednie, odbiorcę, drukarkę oraz kłusownika.


**Nadawca** -> Wysyła liczbę pakietów podaną przy uruchamianiu programu.

**Wierzchołek pośredni** -> Odbiera pakietu z odpowiedniego kanału, a następnie przesyła go do wylosowanego wierzchołka. Zamiast pakietu, może odebrać pułapkę, nadaną przez kłusownika   i zostanie ona tam dopóki do wierzchołka nie nadejdzie kolejny pakiet.<br />

**Odbiorca** -> Odbiera pakiety.

**Skróty/Dłuższe trasy** -> Uruchamiając program należy podać parametry mówiące ile dodatkowych połączeń ma zostać utworzonych. Skrót skróci drogę do odbiorcy, zaś dłuższa trasa     ją wydłuży.<br />

**Drukarka** -> Drukuje komunikaty nadawane przez wierzchołki.<br />

**Kłusownik** -> Wątek kłusownika zastawia pułapki w wierzchołkach, które likwidują kolejne pakiety, jakie odbierze wierzchołek.<br />

**Pakiet** -> Ma przebyć drogę od nadawcy do odbiorcy. Każdy pakiet ma taką samą żywotność, podawaną jako parametr przy uruchamianiu programu. Za każdym razem, gdy pakiet dotrze     do wierzhołka, żywotność zmniejsza się o jeden. Gdy jej  wartość będzie równa zero, a pakiet nie dotrze do odbiorcy, pakiet umiera.<br />

Program należy uruchomić z parametrami 𝑛, 𝑑, 𝑏, 𝑘, ℎ.<br />
𝑛 - liczba wierzchołków<br />
𝑑 - liczba skrótów<br />
𝑏 - liczba dłuższych tras<br />
𝑘 - liczba pakietów wysyłanych<br />
ℎ - żywotność pakietu<br />
Przykładowe wywołanie:<br />
go run symulacja.go 12 3 2 7 13
