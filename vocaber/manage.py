from flask_script import Server, Manager 
from main import app, db, VocabItem

manager = Manager(app)
manager.add_command("server", Server())

@manager.shell
def make_shell_context():
    return dict(app=app, db=db, VocabItem=VocabItem)

if __name__ == "__main__":
    manager.run()