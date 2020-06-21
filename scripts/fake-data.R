library(charlatan)

try(setwd(dirname(rstudioapi::getActiveDocumentContext()$path)))
setwd(system("git rev-parse --show-toplevel", intern =T))

df <- ch_generate("name", "job", "currency", n = 100)
testVector <- MissingDataProvider$new()
df$gene <- ch_gene_sequence(n = 100)
df$card <- testVector$make_missing(x = ch_credit_card_number(n = 100))

df$gene <- testVector$make_missing(x = df$gene)

df$longitude <- ch_lon(n=100)
df$latitude <- ch_lat(n=100)
df$n <- ch_double(n=100)

data.table::fwrite(df, "data/data-01.tsv", sep = '\t')