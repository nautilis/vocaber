import datetime
from config import DevConfig
from flask import Flask, request, jsonify
from flask_sqlalchemy import SQLAlchemy 
from sqlalchemy import desc, func, extract, and_
from datetime import timedelta
import random
from flask_cors import CORS

app = Flask(__name__)
app.config.from_object(DevConfig)
db = SQLAlchemy(app)
CORS(app)

class VocabItem(db.Model):
    id = db.Column(db.Integer(), primary_key=True)
    value = db.Column(db.String(255))
    created = db.Column(db.DateTime(), default=datetime.datetime.now)
    knownit = db.Column(db.Integer(), default=0)

    def __init__(self, value):
        self.value = value

    def __repr___(self):
        return "<VocabItem `{}`".format(self.value)

    def to_dict(self):
        return {c.name: getattr(self, c.name, None) for c in self.__table__.columns}

    @classmethod
    def save(cls,item):
        db.session.add(item)
        db.session.commit()

    @classmethod
    def count(cls):
        today = datetime.datetime.today()
        year = today.year
        month = today.month
        day = today.day
        return cls.query.filter(and_(
            extract('year', cls.created) == year,
            extract('month', cls.created) == month,
            extract('day', cls.created) == day
        )).count() 

    @classmethod
    def select_by_time(cls, year, month, day):
        items = cls.query.filter(and_(
            extract('year', cls.created) == year,
            extract('month', cls.created) == month,
            extract('day', cls.created) == day 
        )).order_by(desc(cls.created)).all()
        return items

    @classmethod
    def known(cls, item_id):
        cls.query.filter_by(id=item_id).update({
            "knownit": cls.knownit + 1 
        })
        db.session.commit()

    @classmethod
    def get_not_master(cls):
        return cls.query.filter(cls.knownit < 10).all()

    @classmethod
    def delete_item(cls, item_id):
        item = cls.query.filter_by(id=item_id).first()
        if(item):
            db.session.delete(item)
            db.session.commit()


@app.route('/item', methods=['POST','GET'])
def hello_world():
    try:
        item = request.form.get("item")
        token = request.form.get("token")
        print(item, "<=======>", token)
        if(token != DevConfig.TOKEN):
            res = {"result": 0 }
            return jsonify(res)
        else:
            new_item = VocabItem(item)
            VocabItem.save(new_item)
            count = VocabItem.count()
            res = {"result": count}
    except Exception:
        raise
        res = {"result": 0}

    return jsonify(res)

@app.route('/items_by_subday')
def get_yesterday_items():
    days = int(request.args.get("subday"))
    today = datetime.datetime.today() 
    yesterday = today - timedelta(days=days)
    year = yesterday.year
    month = yesterday.month
    day = yesterday.day
    items = VocabItem.select_by_time(year, month, day) 
    resp = {}
    resp["items"] = []
    for item in items:
        resp["items"].append(item.to_dict())
    return jsonify(resp) 



@app.route('/known_it', methods=["POST"])
def known_it():
    token = request.form.get("token")
    id = request.form.get("itemid")
    if(token != DevConfig.TOKEN):
        res = {"result": "failed"}
    else:
        print(id, "<====>", token)
        VocabItem.known(id)
        res = {"result": "success"}
    return jsonify(res)

@app.route('/get_not_master')
def get_not_master():
    items = VocabItem.get_not_master()
    random.shuffle(items)
    res_items = items[0:20]
    res = {}
    res["items"] = []
    for item in res_items:
        res["items"].append(item.to_dict())

    return jsonify(res)

@app.route('/get_today_count')
def get_today_count():
    count = VocabItem.count()
    result = {"result": count}
    return jsonify(result)

@app.route('/delete_item', methods=["POST"])
def delete_item():
    token = request.form.get("token")
    id = request.form.get("itemid")
    if(token != DevConfig.TOKEN):
        res = {"result": "failed"}
    else:
        VocabItem.delete_item(id)
        res = {"result": "success"} 
    return jsonify(res)