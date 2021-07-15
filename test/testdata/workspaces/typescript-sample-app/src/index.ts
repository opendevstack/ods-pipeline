var express = require("express")

const app = express();

app.get("/", (req: any, res: any) => {
    res.send("Hello world!");
});

app.listen(8080, () => {
    console.log(`Server is listening on port 8080!`);
});
