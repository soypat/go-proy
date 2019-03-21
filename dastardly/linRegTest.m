% LinRegTest https://stattrek.com/matrix-algebra/covariance-matrix.aspx
randVar = [90 60 90
    90 90 30
    60 60 60
    60 60 90
    30 30 30];
[n, D] = size(randVar);
A = randVar;
onesColumn = ones(n,1); %Vector columna de unos
%5 students, 3 tests
%% First we transfor raw scores to deviation scores
deviationScores = A - onesColumn*onesColumn.'*A/n;
DeviationX = A - sum(A)/n;

%% Now we calculate the deviation score sums of squares matrix

deviationScores.'*deviationScores

%% and then the Variance Covariance Matrix:


COV = deviationScores.'*deviationScores/n
getCovariance(A.')

function [covariance] = getCovariance(X)
    % X (D,n)
    [D, n]=size(X);
    X=X.';
%     Iv = ones(n,1);
    DeviationX = X - sum(X)/n;
    covariance = DeviationX.'*DeviationX/n;
    return
end

