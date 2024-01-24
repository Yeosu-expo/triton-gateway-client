package setting

/* ----- Server Setting ----- */
const ServerPort string = "80"

const ModelPath string = "./models/model_list.json"

/* ----- Triton Server Setting ----- */
const TritonUrl string = "localhost:2000"

const batchSize int = 1
const Samples int = 1
const Steps int = 45
const GuidanceScale float64 = 7.5
