unit ufrmMain;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, FileUtil, Forms, Controls, Graphics, Dialogs, StdCtrls,
  ExtCtrls, CheckLst, Menus;

type

  { TFrmMain }

  TFrmMain = class(TForm)
    BtnGen: TButton;
    BtnUpdate: TButton;
    BtnDelete: TButton;
    BtnGetOne: TButton;
    BtnGetAll: TButton;
    BtnInsert: TButton;
    ClbFields: TCheckListBox;
    MmiUnSelectAll: TMenuItem;
    MmiSelectAll: TMenuItem;
    MmoResult: TMemo;
    MmoStruct: TMemo;
    MmoSQL: TMemo;
    Panel1: TPanel;
    Pn2: TPanel;
    Pn1: TPanel;
    Pn3: TPanel;
    Pn4: TPanel;
    PopFields: TPopupMenu;
    procedure BtnDeleteClick(Sender: TObject);
    procedure BtnGenClick(Sender: TObject);
    procedure BtnGetAllClick(Sender: TObject);
    procedure BtnGetOneClick(Sender: TObject);
    procedure BtnInsertClick(Sender: TObject);
    procedure BtnUpdateClick(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure MmiSelectAllClick(Sender: TObject);
    procedure MmiUnSelectAllClick(Sender: TObject);
  private

  public

  end;

var
  FrmMain: TFrmMain;

implementation

{$R *.lfm}

{ TFrmMain }

procedure TFrmMain.BtnGenClick(Sender: TObject);
begin
  // BtnGenClick
end;

procedure TFrmMain.BtnGetAllClick(Sender: TObject);
begin
  // BtnGetAllClick
end;

procedure TFrmMain.BtnDeleteClick(Sender: TObject);
begin
  // BtnDeleteClick
end;

procedure TFrmMain.BtnGetOneClick(Sender: TObject);
begin
  // BtnGetOneClick
end;

procedure TFrmMain.BtnInsertClick(Sender: TObject);
begin
  // BtnInsertClick
end;

procedure TFrmMain.BtnUpdateClick(Sender: TObject);
begin
  // BtnUpdateClick
end;

procedure TFrmMain.FormCreate(Sender: TObject);
begin
  // FormCreate
end;

procedure TFrmMain.MmiSelectAllClick(Sender: TObject);
begin
  // MmiSelectAllClick
end;

procedure TFrmMain.MmiUnSelectAllClick(Sender: TObject);
begin
  // MmiUnSelectAllClick
end;

end.

