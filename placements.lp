#show placed/2.

piece(s1;s2;s3;s4;h1;v1;v2;v3;v4;big).

xv(1;2;3;4).
yv(1;2;3;4;5).

loc(p(X,Y)):- xv(X), yv(Y).


type(s1,small).
type(s2,small).
type(s3,small).
type(s4,small).

type(h1,horiz).

type(v1,vert).
type(v2,vert).
type(v3,vert).
type(v4,vert).

type(big,big).

full(p(X,Y),P):- placed(P,p(X,Y)),type(P,small).

full(p(X,Y),P):- placed(P,p(X,Y)),type(P,horiz).
full(p(X,Y),P):- placed(P,p(X+1,Y)),type(P,horiz).


full(p(X,Y),P):- placed(P,p(X,Y)),type(P,vert).
full(p(X,Y),P):- placed(P,p(X,Y+1)),type(P,vert).


full(p(X,Y),P):- placed(P,p(X,Y)),type(P,big).
full(p(X,Y),P):- placed(P,p(X,Y+1)),type(P,big).
full(p(X,Y),P):- placed(P,p(X+1,Y)),type(P,big).
full(p(X,Y),P):- placed(P,p(X+1,Y+1)),type(P,big).

:- full(p(X,_),_), not xv(X).
:- full(p(_,Y),_), not yv(Y).

canplace(P,p(X,Y)) :- type(P,small),xv(X),yv(Y).
canplace(P,p(X,Y)) :- type(P,horiz), xv(X),yv(Y), X < 4.
canplace(P,p(X,Y)) :- type(P,vert),xv(X),yv(Y), Y < 5.
canplace(P,p(X,Y)) :- type(P,big), xv(X),yv(Y), X < 4, Y < 5.


1{placed(P,POS) :  canplace(P,POS)}1 :- piece(P).

:- full(POS,P1), full(POS,P2), P1 != P2.